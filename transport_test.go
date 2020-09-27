package tenuki_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/podhmo/tenuki"
)

func TestHandlerRoundTripper(t *testing.T) {
	f := tenuki.New(t)
	f.Client = &http.Client{
		Transport: &tenuki.HandlerTripper{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "hello world"+r.URL.Query().Get("suffix"))
			}),
		},
	}

	req := f.NewRequest("GET", "http:", nil)
	q := req.URL.Query()
	q.Add("suffix", " !!")
	req.URL.RawQuery = q.Encode()

	res := f.Do(req)

	want := "hello world !!"
	got := string(f.Extract().Bytes(res))
	if want != got {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
