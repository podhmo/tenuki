package capture

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/podhmo/tenuki/capture/style"
)

// for text output
type bytesInfo struct {
	req *http.Request
	b   []byte
}

type bytesState struct {
	req *http.Request
	b   []byte
}

func (s *bytesState) Info() style.Info {
	return s
}
func (s *bytesState) Merge(res style.Info) style.Info {
	log.Println("bytesState.Merge is not implemented")
	return s
}

func (s *bytesState) Encode() ([]byte, error) {
	return s.b, nil
}

func (s *bytesState) Emit(w io.Writer) error {
	if _, err := w.Write(s.b); err != nil {
		return err
	}
	return nil
}

func (s *bytesState) EmitBoth(w io.Writer, res style.State) error {
	// acting as commit like function, so emitting all states.
	if err := s.Emit(w); err != nil {
		return fmt.Errorf("prev emit %w", err)
	}
	fmt.Fprint(w, "\n----------------------------------------\n\n")
	if err := res.Emit(w); err != nil {
		return fmt.Errorf("this emit %w", err)
	}
	return nil
}

func (s *bytesState) HandleError(open func() (io.WriteCloser, error), err error) {
	f, openErr := open()
	if openErr != nil {
		log.Printf("something wrong, when open file %+v", err)
		return
	}

	if err != nil {
		defer f.Close()
		s.dumpHeader(f)
		fmt.Fprintf(f, "%+v\n", err)
	}
}

func (s *bytesState) dumpHeader(w io.Writer) {
	req := s.req

	reqURI := req.RequestURI
	if reqURI == "" {
		reqURI = req.URL.RequestURI()
	}
	method := req.Method
	if method == "" {
		method = "GET"
	}
	fmt.Fprintf(w, "%s %s HTTP/%d.%d\r\n", method,
		reqURI, req.ProtoMajor, req.ProtoMinor)
}

// for json output
type jsonState struct {
	info style.Info
}

func (s *jsonState) Encode() ([]byte, error) {
	return s.encodeInfo(s.info)
}
func (s *jsonState) encodeInfo(info style.Info) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(info); err != nil {
		return nil, fmt.Errorf("encode json, %w", err)
	}
	return b.Bytes(), nil
}
func (s *jsonState) Emit(f io.Writer) error {
	info := s.info
	b, err := s.encodeInfo(info)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("write json, %w", err)
	}
	return nil
}
func (s *jsonState) EmitBoth(f io.Writer, res style.State) error {
	info := s.info.Merge(res.Info())
	b, err := s.encodeInfo(info)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("write json, %w", err)
	}
	return nil
}

func (s *jsonState) Info() style.Info {
	return s.info
}
