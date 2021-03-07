package capture

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/podhmo/tenuki/capture/style"
)

// for text output
type keepPrevState struct {
	prev style.State
	this style.State
}

func (s *keepPrevState) Encode() ([]byte, error) {
	// the Assumption that keep's Encode() is already called.
	return s.this.Encode()
}
func (s *keepPrevState) Info() style.Info {
	return s.this.Info()
}
func (s *keepPrevState) Emit(w io.Writer) error {
	// acting as commit like function, so emitting all states.
	if err := s.prev.Emit(w); err != nil {
		return fmt.Errorf("prev emit %w", err)
	}
	fmt.Fprint(w, "\n----------------------------------------\n\n")
	if err := s.this.Emit(w); err != nil {
		return fmt.Errorf("this emit %w", err)
	}
	return nil
}

type bytesState struct {
	req *http.Request
	b   []byte
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
func (s *bytesState) Info() style.Info {
	return s
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
