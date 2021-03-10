package main

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/podhmo/tenuki"
	"github.com/podhmo/tenuki/difftest"
)

// func init() {
// 	tenuki.DefaultLayout = capture.OpenAPILayout
// }

func Test200(t *testing.T) {
	targetHandler := Handler200

	f := tenuki.New(t)
	req := f.NewRequest("GET", "", nil)
	res := f.DoHandlerFunc(targetHandler, req,
		tenuki.AssertStatus(200),
	)

	want := `
{
    "status" : "success",
    "data" : {
        "posts" : [
            { "id" : 1, "title" : "A blog post", "body" : "Some useful content" },
            { "id" : 2, "title" : "Another blog post", "body" : "More content" }
        ]
     }
}`
	got := f.Extract().JSON(res)
	difftest.AssertGotAndWantString(t, got, want)
}

func Test400(t *testing.T) {
	targetHandler := Handler400

	f := tenuki.New(t)
	req := f.NewJSONRequest("POST", "/articles", strings.NewReader(`{"content": "Some useful content"}`))
	res := f.DoHandlerFunc(targetHandler, req,
		tenuki.AssertStatus(400),
	)

	want := `
{
    "status" : "fail",
    "message" : "bad request",
    "data" : {
        "title": "A title is required"
     }
}`
	got := f.Extract().JSON(res)
	difftest.AssertGotAndWantString(t, got, want)
}

func Test500(t *testing.T) {
	targetHandler := Handler500

	f := tenuki.New(t)
	req := f.NewRequest("GET", "", nil)
	res := f.DoHandlerFunc(targetHandler, req,
		tenuki.AssertStatus(500),
	)

	want := `
{
    "status" : "error",
    "message" : "Unable to communicate with database"
}`
	got := f.Extract().JSON(res)
	difftest.AssertGotAndWantString(t, got, want)
}

func TestNetworkUnreached(t *testing.T) {
	ts := httptest.NewServer(tenuki.NewCloseConnectionHandler(t))
	defer ts.Close()

	f := tenuki.New(t)
	req := f.NewRequest("GET", ts.URL, nil)
	f.Do(req,
		tenuki.AssertError(func(err error) error {
			if err == nil {
				return fmt.Errorf("something wrong !!!!!!!!!!!!")
			}
			t.Logf("ok error is occured %q", err)
			return nil
		}),
	)
}
