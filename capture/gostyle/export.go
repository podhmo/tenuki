package gostyle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func ExtractRequestInfo(req *http.Request, body io.ReadCloser) (Info, error) {
	info := InfoFromInterface(req, []string{
		"URL", "Body", "GetBody", "Close", "Trailer", "TLS", "Cancel", "Response", "ctx",
	})

	if body != nil {
		ct := strings.ToLower(req.Header.Get("Content-Type"))
		body, err := parseBody(body, ct)
		if err != nil {
			return info, fmt.Errorf("parse body, %w", err)
		}
		if body != nil {
			info["Body"] = body
		}
	}
	return info, nil
}

func ExtractResponseInfo(resp *http.Response, body io.ReadCloser) (Info, error) {
	info := InfoFromInterface(resp, []string{
		"Close", "Body", "Trailer", "Request", "TLS", "Request",
	})

	if body != nil {
		ct := strings.ToLower(resp.Header.Get("Content-Type"))
		body, err := parseBody(body, ct)
		if err != nil {
			return info, fmt.Errorf("parse body, %w", err)
		}
		if body != nil {
			info["Body"] = body
		}
	}
	return info, nil
}

func parseBody(body io.ReadCloser, contentType string) (interface{}, error) {
	ct := strings.SplitN(contentType, "+", 2)[0]
	switch ct {
	case "application/json", "text/json":
		var ob interface{}
		if err := json.NewDecoder(body).Decode(&ob); err != nil {
			return nil, fmt.Errorf("unmarshal json body, %w", err)
		}
		return ob, nil
	default:
		if strings.HasPrefix(ct, "text/") || ct == "application/x-www-form-urlencoded" {
			var b bytes.Buffer
			if _, err := b.ReadFrom(body); err != nil {
				return nil, fmt.Errorf("read body, %w", err)
			}
			// TODO: truncate?
			return b.String(), nil
		} else {
			log.Printf("Content-Type=%q is not supported, so just ignored", ct)
			return nil, nil
		}
	}
}
