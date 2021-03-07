package tenuki_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/podhmo/tenuki"
)

func TestHandlerRoundTripper(t *testing.T) {
	transport := tenuki.HandlerTripperFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world"+r.URL.Query().Get("suffix"))
	})

	f := tenuki.New(t)
	req := f.NewRequest("GET", "", nil)

	q := req.URL.Query()
	q.Add("suffix", " !!")
	req.URL.RawQuery = q.Encode()

	f.Client = &http.Client{Transport: transport}
	res := f.Do(req)

	want := "hello world !!"
	got := string(f.Extract().Bytes(res))
	if want != got {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
