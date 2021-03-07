package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/podhmo/tenuki"
	"github.com/podhmo/tenuki/difftest"
)

func Test(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		tenuki.Render(w, r).JSON(200, map[string]string{"message": "hello world"})
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))

	f := tenuki.New(t)
	req := f.NewRequest("GET", ts.URL+"/hello", nil)
	res := f.Do(req,
		tenuki.AssertStatus(200),
	)

	var got map[string]string
	f.Extract().JSON(res, &got)

	difftest.AssertGotAndWantString(t,
		got,
		`{"message": "hello world"}`,
	)
}
