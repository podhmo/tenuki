package reqtest_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/podhmo/reqtest"
)

func TestHandlerRoundTripper(t *testing.T) {
	f := reqtest.New(t)
	f.Client = &http.Client{
		Transport: &reqtest.HandlerTripper{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "hello world"+r.URL.Query().Get("suffix"))
			}),
		},
	}

	req := f.NewRequest(
		"GET", "http:", nil,
		reqtest.WithQuery(func(q url.Values) {
			q.Add("suffix", " !!")
		}),
	)

	res := f.Do(req)

	want := "hello world !!"
	got := string(f.Extract().Bytes(res))
	if want != got {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
