package openapistyle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"sort"
	"strings"
)

// TODO: trim security information
func toPaths(req *http.Request, info *Info) (Paths, error) {
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

		info.Method = valueOrDefault(req.Method, "GET")
		info.HTTPVersion = req.Proto
		info.HeaderSize = -1 // TODO
		info.BodySize = -1   // TODO
	}

	r := Paths{}
	pathItem, err := toPathItem(req, info)
	if err != nil {
		return r, fmt.Errorf("extract pathItem, %w", err)
	}
	r[req.URL.Path] = pathItem
	return r, nil
}
func toPathItem(req *http.Request, info *Info) (PathItem, error) {
	r := PathItem{}

	op, err := toOperation(req, info)
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

func toOperation(req *http.Request, info *Info) (Operation, error) {
	r := Operation{}

	if req.Body != nil {
		ct, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			log.Printf("parse content type, %+v", err)
			ct = strings.ToLower(req.Header.Get("Content-Type"))
		}
		info.ContentType = ct

		content, err := toContent(req, info)
		if err != nil {
			return r, fmt.Errorf("extract content, %w", err)
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

func toContent(req *http.Request, info *Info) (MediaType, error) {
	r := MediaType{}
	body := req.Body
	switch ct := info.ContentType; ct {
	case "application/json", "text/json":
		if err := json.NewDecoder(body).Decode(&r.Example); err != nil {
			return r, fmt.Errorf("unmarshal json body, %w", err)
		}
		return r, nil
	default:
		if strings.HasPrefix(ct, "text/") || ct == "application/x-www-form-urlencoded" {
			var b bytes.Buffer
			if _, err := b.ReadFrom(body); err != nil {
				return r, fmt.Errorf("read body, %w", err)
			}
			r.Example = b.String()
			return r, nil
		}
		log.Printf("Content-Type=%q is not supported, so just ignored", ct)
		return r, nil
	}
}

// func toOperation(req *http.Request, info *Info) (Operation, error) {
// }
