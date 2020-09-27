package reqtest_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/podhmo/reqtest"
)

func TestIt(t *testing.T) {
	type Body struct {
		Message string `json:"message"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqtest.Render(w, r).JSON(Body{Message: "hello world"})
	}))
	defer ts.Close()

	f := reqtest.New(t)
	res := f.Do(f.NewRequest("GET", ts.URL, nil))

	// assertion
	{
		if want, got := http.StatusOK, res.StatusCode; want != got {
			t.Errorf("status code:\nwant\n\t%+v\nbut\n\t%+v", want, got)
		}

		want := Body{Message: "hello world"}
		var got Body
		f.Extract().JSON(res, &got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
		}
	}

	// extract multiple times is also ok.
	{
		want := Body{Message: "hello world"}
		var got Body
		f.Extract().JSON(res, &got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
		}
	}
}
