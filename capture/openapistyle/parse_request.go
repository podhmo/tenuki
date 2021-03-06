package openapistyle

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// TODO: trim security information
func toPaths(req *http.Request, body io.ReadCloser) (Paths, error) {
	r := Paths{}
	pathItem, err := toPathItem(req, body)
	if err != nil {
		return r, fmt.Errorf("extract pathItem, %w", err)
	}
	r[req.URL.Path] = pathItem
	return r, nil
}
func toPathItem(req *http.Request, body io.ReadCloser) (PathItem, error) {
	r := PathItem{}

	op, err := toOperation(req, body)
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

func toOperation(req *http.Request, body io.ReadCloser) (Operation, error) {
	r := Operation{}

	if body != nil {
		content, err := toContent(req, body)
		if err != nil {
			return r, fmt.Errorf("extract content, %w", err)
		}
		r.RequestBody = &RequestBody{
			Content: map[string]MediaType{
				req.Header.Get("Content-Type"): content,
			},
		}
	}

	// query, header, (path), cookie
	if q := req.URL.Query(); q != nil {
		for k, vs := range q {
			examples := make([]interface{}, len(vs))
			for i, v := range vs {
				examples[i] = v
			}
			r.Parameters = append(r.Parameters, Parameter{
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
			r.Parameters = append(r.Parameters, Parameter{
				Name:     k,
				In:       "header",
				Examples: examples,
			})
		}
	}
	if cookies := req.Cookies(); len(cookies) > 0 {
		for _, cookie := range cookies {
			r.Parameters = append(r.Parameters, Parameter{
				Name:     cookie.Name,
				In:       "cookie",
				Examples: []interface{}{cookie.Raw}, // invalid?
			})
		}
	}
	return r, nil
}

func toContent(req *http.Request, body io.ReadCloser) (MediaType, error) {
	r := MediaType{}
	return r, nil
}

// func toOperation(req *http.Request, body io.ReadCloser) (Operation, error) {
// }