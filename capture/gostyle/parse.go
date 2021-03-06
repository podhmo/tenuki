package gostyle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// reflection version

type Info map[string]interface{}

// for interface
func (i Info) Info() interface{} {
	return nil
}

func parseRequest(req *http.Request, body io.Reader) (Info, error) {
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

func parseResponse(resp *http.Response, body io.Reader) (Info, error) {
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

func InfoFromInterface(ptr interface{}, excludes []string) Info {
	rt := reflect.TypeOf(ptr).Elem()
	rv := reflect.ValueOf(ptr).Elem()
	info := Info{}

toplevel:
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		for _, name := range excludes {
			if name == rf.Name {
				continue toplevel
			}
		}
		info[rf.Name] = rv.Field(i).Interface()
	}
	return info
}

func parseBody(body io.Reader, contentType string) (interface{}, error) {
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
