package reqtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type ExtractFacade struct {
	T     *testing.T
	cache map[string][]byte
}

func (f *Facade) Extract() *ExtractFacade {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.extractor != nil {
		return f.extractor
	}
	f.extractor = &ExtractFacade{
		T:     f.T,
		cache: map[string][]byte{},
	}
	return f.extractor
}

func (f *ExtractFacade) buffer(res *http.Response) io.Reader {
	k := fmt.Sprintf("%p", res) // xxx
	cache := f.cache[k]
	if cache != nil {
		return bytes.NewReader(cache)
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

	cache = b.Bytes()
	f.cache[k] = cache
	return bytes.NewReader(cache)
}

func (f *ExtractFacade) JSON(res *http.Response, ob interface{}) {
	t := f.T
	t.Helper()

	decoder := json.NewDecoder(f.buffer(res))
	if err := decoder.Decode(&ob); err != nil {
		t.Fatalf("!! DecodeJSON: %+v", err)
	}
}
