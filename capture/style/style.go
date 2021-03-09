package style

import (
	"io"
)

type Info interface {
	Merge(res Info) Info

	HandleError(open func() (io.WriteCloser, error), err error)
	// HandleResponse(open func() (io.WriteCloser, error), res *http.Response)
}

type State interface {
	Info() Info

	Encode() ([]byte, error)

	Emit(w io.Writer) error
	EmitBoth(w io.Writer, res State) error
}
