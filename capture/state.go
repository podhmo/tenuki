package capture

import (
	"io"
	"net/http"
)

type State interface {
	Encode() ([]byte, error)
	Emit(w io.WriteCloser) error

	// HandleError(open func() (io.WriteCloser, error), err error) error
	// HandleResponse(w io.Writer, res *http.Response) error
}

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

// func (l bytesLazy) HandleError(open func() (io.Writer, error), err error) error {
// 	wt.dumpHeader(f, req)
// 	fmt.Fprintf(f, "%+v\n", err)
// 	return nil
// }
