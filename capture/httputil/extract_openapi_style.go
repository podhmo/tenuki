package httputil

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/podhmo/tenuki/capture/openapi"
)

// TODO: trim security information

func extractOpenAPIPaths(req *http.Request, body io.ReadCloser) (openapi.Paths, error) {
	r := openapi.Paths{}
	pathItem, err := extractOpenAPIPathItem(req, body)
	if err != nil {
		return r, fmt.Errorf("extract pathItem, %w", err)
	}
	r[req.URL.Path] = pathItem
	return r, nil
}
func extractOpenAPIPathItem(req *http.Request, body io.ReadCloser) (openapi.PathItem, error) {
	r := openapi.PathItem{}

	op, err := extractOpenAPIOperation(req, body)
	if err != nil {
		return r, fmt.Errorf("extract operation, %w", err)
	}

	switch method := strings.ToLower(req.Method); method {
	case "get":
		r.Get = &op
	case "post":
		r.Post = &op
	case "delete":
		r.Delete = &op
	case "options":
		r.Options = &op
	case "head":
		r.Head = &op
	case "patch":
		r.Patch = &op
	case "trace":
		r.Trace = &op
	default:
		log.Printf("unknown method %s, treated as GET", method)
		r.Get = &op
	}
	return r, nil
}

func extractOpenAPIOperation(req *http.Request, body io.ReadCloser) (openapi.Operation, error) {
	r := openapi.Operation{}
	content, err := extractOpenAPIContent(req, body)
	if err != nil {
		return r, fmt.Errorf("extract content, %w", err)
	}
	r.RequestBody = &openapi.RequestBody{
		Content: map[string]openapi.MediaType{
			req.Header.Get("Content-Type"): content,
		},
	}

	// query, header, (path), cookie
	if q := req.URL.Query(); q != nil {
		for k, vs := range q {
			examples := make([]interface{}, len(vs))
			for i, v := range vs {
				examples[i] = v
			}
			r.Parameters = append(r.Parameters, openapi.Parameter{
				Name:     k,
				In:       "query",
				Examples: examples,
			})
		}
	}
	if len(req.Header) > 0 {
		for k, vs := range req.Header {
			examples := make([]interface{}, len(vs))
			for i, v := range vs {
				examples[i] = v
			}
			r.Parameters = append(r.Parameters, openapi.Parameter{
				Name:     k,
				In:       "header",
				Examples: examples,
			})
		}
	}
	if cookies := req.Cookies(); len(cookies) > 0 {
		for _, cookie := range cookies {
			r.Parameters = append(r.Parameters, openapi.Parameter{
				Name:     cookie.Name,
				In:       "cookie",
				Examples: []interface{}{cookie.Raw}, // invalid?
			})
		}
	}
	return r, nil
}

func extractOpenAPIContent(req *http.Request, body io.ReadCloser) (openapi.MediaType, error) {
	r := openapi.MediaType{}
	return r, nil
}

// func extractOpenAPIOperation(req *http.Request, body io.ReadCloser) (openapi.Operation, error) {
// }
