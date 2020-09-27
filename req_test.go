package reqtest_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/podhmo/reqtest"
)

func TestIt(t *testing.T) {
	type body struct {
		Message string `json:"message"`
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		reqtest.Render(w, r).JSON(body{Message: "hello world"})
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	f := reqtest.New(t)
	req := f.NewRequest("GET", ts.URL, nil)
	res := f.Do(req)

	if want, got := http.StatusOK, res.StatusCode; want != got {
		t.Errorf("status code:\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}

	want := body{Message: "hello world"}
	var got body
	f.Extract().JSON(res, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
