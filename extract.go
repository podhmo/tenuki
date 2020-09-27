package reqtest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type ExtractFacade struct {
	T     *testing.T
	cache map[*http.Response]*bytes.Reader
}

func (f *Facade) Extract() *ExtractFacade {
	return &ExtractFacade{
		T:     f.T,
		cache: map[*http.Response]*bytes.Reader{},
	}
}

func (f *ExtractFacade) buffer(res *http.Response) io.Reader {
	cache := f.cache[res]
	if cache != nil {
		cache.Seek(0, 0)
		return cache
	}

	t := f.T
	t.Helper()

	var b bytes.Buffer
	if _, err := io.Copy(&b, res.Body); err != nil {
		t.Fatalf("!! buffer: %+v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("!! DecodeJSON Close: %+v", err)
		}
	}()

	cache = bytes.NewReader(b.Bytes())
	f.cache[res] = cache
	return cache
}

func (f *ExtractFacade) JSON(res *http.Response, ob interface{}) {
	t := f.T
	t.Helper()

	decoder := json.NewDecoder(f.buffer(res))
	if err := decoder.Decode(&ob); err != nil {
		t.Fatalf("!! DecodeJSON: %+v", err)
	}
}
