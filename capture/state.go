package capture

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/podhmo/tenuki/capture/style"
)

// for text output
type bytesState struct {
	req *http.Request
	b   []byte
}

func (s *bytesState) Encode() ([]byte, error) {
	return s.b, nil
}
func (s *bytesState) Emit(w io.WriteCloser) error {
	if _, err := w.Write(s.b); err != nil {
		return err
	}
	defer w.Close()
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
