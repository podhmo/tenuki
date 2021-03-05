package tenuki_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/podhmo/tenuki"
)

func TestIt(t *testing.T) {
	type Body struct {
		Message string `json:"message"`
	}

	sumHandler := func(w http.ResponseWriter, r *http.Request) {
		var xs []int
		tenuki.DecodeJSON(r.Body, &xs)
		n := 0
		for i := range xs {
			n += xs[i]
		}
		tenuki.Render(w, r).JSON(200, Body{Message: fmt.Sprintf("sum is %d", n)})
	}
	ts := httptest.NewServer(http.HandlerFunc(sumHandler))
	defer ts.Close()

	f := tenuki.New(t)
	res := f.Do(
		f.NewRequest("Post", ts.URL, strings.NewReader(`[1,2,3]`)),
		tenuki.AssertStatus(http.StatusOK),
	)

	// assertion
	want := Body{Message: "sum is 6"}
	var got Body
	f.Extract().JSON(res, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}

	// extract multiple times is also ok.
	{
		want := Body{Message: "sum is 6"}
		var got Body
		f.Extract().JSON(res, &got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
		}
	}
}
