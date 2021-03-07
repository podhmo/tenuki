package capture

import (
	"io"
)

type State interface {
	Encode() ([]byte, error)
	Emit(w io.WriteCloser) error
}

// for text output
type bytesLazy []byte

func (l bytesLazy) Encode() ([]byte, error) {
	return l, nil
}
func (l bytesLazy) Emit(w io.WriteCloser) error {
	if _, err := w.Write(l); err != nil {
		return err
	}
	defer w.Close()
	return nil
}
