package tenuki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

type ExtractFacade struct {
	T *testing.T

	cache map[string][]byte
	mu    sync.Mutex
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

func (f *ExtractFacade) buffer(res *http.Response) *bytes.Buffer {
	f.mu.Lock()
	defer f.mu.Unlock()

	k := fmt.Sprintf("%p", res) // xxx
	cache := f.cache[k]
	if cache != nil {
		return bytes.NewBuffer(cache)
	}

	t := f.T
	t.Helper()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		t.Fatalf("!! buffer: %+v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("!! DecodeJSON Close: %+v", err)
		}
	}()

	cache = buf.Bytes()
	f.cache[k] = cache
	return bytes.NewBuffer(cache)
}

func (f *ExtractFacade) BindJSON(res *http.Response, ob interface{}) {
	t := f.T
	t.Helper()

	decoder := json.NewDecoder(f.buffer(res))
	if err := decoder.Decode(&ob); err != nil {
		t.Fatalf("!! DecodeJSON: %+v", err)
	}
}
func (f *ExtractFacade) JSON(res *http.Response) interface{} {
	t := f.T
	t.Helper()
	var ob interface{}
	f.BindJSON(res, &ob)
	return ob
}
func (f *ExtractFacade) Bytes(res *http.Response) []byte {
	t := f.T
	t.Helper()
	return f.buffer(res).Bytes()
}
func (f *ExtractFacade) String(res *http.Response) string {
	t := f.T
	t.Helper()
	return f.buffer(res).String()
}

func DecodeJSON(r io.Reader, ob interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(ob)
}
