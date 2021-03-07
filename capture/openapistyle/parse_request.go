package openapistyle

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"sort"
	"strings"
)

// TODO: trim security information
func toPaths(req *http.Request, body io.Reader) (Paths, error) {
	r := Paths{}
	pathItem, err := toPathItem(req, body)
	if err != nil {
		return r, fmt.Errorf("extract pathItem, %w", err)
	}
	r[req.URL.Path] = pathItem
	return r, nil
}
func toPathItem(req *http.Request, body io.Reader) (PathItem, error) {
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

func toOperation(req *http.Request, body io.Reader) (Operation, error) {
	r := Operation{}

	if body != nil {
		content, err := toContent(req, body)
		if err != nil {
			return r, fmt.Errorf("extract content, %w", err)
		}
		ct, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			log.Printf("parse content type, %+v", err)
			ct = strings.ToLower(req.Header.Get("Content-Type"))
		}
		r.RequestBody = &RequestBody{
			Content: map[string]MediaType{
				ct: content,
			},
		}
	}

	// query, header, (path), cookie
	if q := req.URL.Query(); q != nil {
		keys := make([]string, 0, len(q))
		for k := range q {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			vs := q[k]
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
		keys := make([]string, 0, len(req.Header))
		for k := range req.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			vs := req.Header[k]
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

func toContent(req *http.Request, body io.Reader) (MediaType, error) {
	r := MediaType{}
	return r, nil
}

// func toOperation(req *http.Request, body io.Reader) (Operation, error) {
// }
