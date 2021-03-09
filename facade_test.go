package tenuki_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/podhmo/tenuki"
)

type Body struct {
	Message string `json:"message"`
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	var xs []int
	tenuki.DecodeJSON(r.Body, &xs)
	n := 0
	for i := range xs {
		n += xs[i]
	}
	tenuki.Render(w, r).JSON(200, Body{Message: fmt.Sprintf("sum is %d", n)})
}

func TestDo(t *testing.T) {
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
	f.Extract().BindJSON(res, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}

	// extract multiple times is also ok.
	{
		want := Body{Message: "sum is 6"}
		var got Body
		f.Extract().BindJSON(res, &got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
		}
	}
}

func TestDoWithBody(t *testing.T) {
	body := `
{
    "id": 1,
    "category": {
        "id": 1,
        "name": ""
    },
    "name": "doggie",
    "photoUrls": [
        ""
    ],
    "tags": [
        {
            "id": 1,
            "name": ""
        }
    ],
    "status": "available"
}`
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, r.Body)
		defer r.Body.Close()
	}
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	f := tenuki.New(t)
	req := f.NewRequest("Post", ts.URL, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	res := f.Do(req,
		tenuki.AssertStatus(http.StatusOK),
	)
	_ = res
}

func TestDoHandler(t *testing.T) {
	f := tenuki.New(t)
	res := f.DoHandlerFunc(sumHandler,
		f.NewRequest("Post", "", strings.NewReader(`[1,2,3]`)),
		tenuki.AssertStatus(http.StatusOK),
	)

	// assertion
	want := Body{Message: "sum is 6"}
	var got Body
	f.Extract().BindJSON(res, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}

	// extract multiple times is also ok.
	{
		want := Body{Message: "sum is 6"}
		var got Body
		f.Extract().BindJSON(res, &got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
		}
	}
}
