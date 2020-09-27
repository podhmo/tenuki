package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/podhmo/tenuki"
)

func Test(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		tenuki.Render(w, r).JSON(200, map[string]string{"message": "hello world"})
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))

	f := tenuki.New(t)
	req := f.NewRequest("GET", ts.URL+"/hello", nil)
	res := f.Do(req)

	if want, got := 200, res.StatusCode; want != got {
		t.Errorf("status code\nwant\n\t%d\nbut\n\t%d", want, got)
	}

	want := map[string]string{"message": "hello world"}
	var got map[string]string
	f.Extract().JSON(res, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
