package style

import (
	"io"
)

type Info interface {
	HandleError(open func() (io.WriteCloser, error), err error)
	// HandleResponse(open func() (io.WriteCloser, error), res *http.Response)
}

type State interface {
	Info() Info

	Encode() ([]byte, error)
	Emit(w io.WriteCloser) error
}
