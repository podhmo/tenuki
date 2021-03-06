package httputil

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/podhmo/tenuki/capture/openapi"
)

type JSONRequestInfo struct {
	Method      string `json:"method"`
	URL         string `json:"url"`
	HTTPVersion string `json:"httpVersion,omitempty"`
	HeaderSize  int64  `json:"headerSize,omitempty" default:"-1"`
	BodySize    int64  `json:"bodySize,omitempty" default:"-1"`

	Paths openapi.Paths `json:"paths"`
}

func ExtractJSONRequestInfo(req *http.Request, body io.ReadCloser) (JSONRequestInfo, error) {
	info := JSONRequestInfo{}

	{
		reqURI := req.RequestURI
		if reqURI == "" {
			reqURI = req.URL.RequestURI()
		}

		absRequestURI := strings.HasPrefix(req.RequestURI, "http://") || strings.HasPrefix(req.RequestURI, "https://")
		if !absRequestURI {
			reqURI = fmt.Sprintf("%s://%s%s", valueOrDefault(req.URL.Scheme, "https"), req.URL.Host, reqURI)
		}
		info.URL = reqURI
	}

	info.Method = valueOrDefault(req.Method, "GET")
	info.HTTPVersion = fmt.Sprintf("HTTP/%d.%d", req.ProtoMajor, req.ProtoMinor)
	paths, err := extractOpenAPIPaths(req, body)
	if err != nil {
		return info, fmt.Errorf("extract paths, %w", err)
	}
	info.Paths = paths
	return info, nil
}

func DumpRequestJSON(req *http.Request, body bool) (JSONRequestInfo, error) {
	var err error
	save := req.Body
	{
		if !body || req.Body == nil {
			req.Body = nil
		} else {
			save, req.Body, err = drainBody(req.Body)
			if err != nil {
				return JSONRequestInfo{}, err
			}
		}
	}
	info, err := ExtractJSONRequestInfo(req, save)
	if err != nil {
		return JSONRequestInfo{}, err
	}
	return info, nil
}

type ResponseInfo struct {
}

func ExtractResponseInfo(resp *http.Response) ResponseInfo {
	info := ResponseInfo{}
	return info
}
func DumpResponseJSON(resp *http.Response, body bool) (ResponseInfo, error) {
	info := ExtractResponseInfo(resp)
	var err error
	save := resp.Body
	savecl := resp.ContentLength

	// TODO: content-type, json の場合は取り出す
	{
		if !body {
			// For content length of zero. Make sure the body is an empty
			// reader, instead of returning error through failureToReadBody{}.
			if resp.ContentLength == 0 {
				resp.Body = emptyBody
			} else {
				resp.Body = failureToReadBody{}
			}
		} else if resp.Body == nil {
			resp.Body = emptyBody
		} else {
			save, resp.Body, err = drainBody(resp.Body)
			if err != nil {
				return info, err
			}
		}

		if err == errNoBody {
			err = nil
		}
		resp.Body = save
		resp.ContentLength = savecl
	}

	if err != nil {
		return info, err
	}
	return info, nil
}
