package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JSONState struct {
	Info interface {
		Info() interface{}
	}
}

func (s *JSONState) Encode() ([]byte, error) {
	info := s.Info
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(info); err != nil {
		return nil, fmt.Errorf("encode json, %w", err)
	}
	return b.Bytes(), nil
}

func (s *JSONState) Emit(f io.WriteCloser) error {
	defer f.Close()
	b, err := s.Encode()
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("write json, %w", err)
	}
	return nil
}

func DumpRequestJSON(
	req *http.Request,
	body bool,
	extractInfo func(*http.Request, io.Reader) (interface{ Info() interface{} }, error),
) (*JSONState, error) {
	var err error
	save := req.Body
	{
		if !body || req.Body == nil {
			req.Body = nil
		} else {
			save, req.Body, err = drainBody(req.Body)
			if err != nil {
				return nil, err
			}
		}
	}
	info, err := extractInfo(req, save)
	if err != nil {
		return nil, fmt.Errorf("extract request info, %w", err)
	}
	return &JSONState{Info: info}, nil
}

func DumpResponseJSON(
	resp *http.Response,
	body bool,
	extractInfo func(*http.Response, io.Reader) (interface{ Info() interface{} }, error),
) (*JSONState, error) {
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
				return nil, err
			}
		}
	}

	info, err := extractInfo(resp, resp.Body)
	resp.Body = save
	resp.ContentLength = savecl
	if err != nil {
		return nil, fmt.Errorf("extract response info, %w", err)
	}
	return &JSONState{Info: info}, nil
}
