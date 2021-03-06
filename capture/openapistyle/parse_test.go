package openapistyle

import (
	"io"
	"net/http"
	"testing"

	"github.com/podhmo/tenuki/difftest"
)

func TestToPaths(t *testing.T) {
	url := "https://example.net"

	t.Run("queryString", func(t *testing.T) {
		path := "/xxx/yyy"
		req, _ := http.NewRequest("GET", url+path, nil)
		q := req.URL.Query()
		q.Add("q", "x")
		q.Add("q", "y")
		q.Add("val", "z")
		q.Add("escaped", `He says "I'm not robot!", but ...`)
		req.URL.RawQuery = q.Encode()
		got := toPathsWithValidation(t, req, nil)
		want := []byte(`
	{
	  "/xxx/yyy": {
	    "get": {
	      "parameters": [
	        {
	          "in": "query",
	          "name": "escaped",
	          "examples": [
	            "He says \"I'm not robot!\", but ..."
	          ]
	        },
	        {
	          "in": "query",
	          "name": "q",
	          "examples": [
	            "x",
	            "y"
	          ]
	        },
	        {
	          "in": "query",
	          "name": "val",
	          "examples": [
	            "z"
	          ]
	        }
	      ]
	    }
	  }
	}
	`)
		difftest.AssertGotAndWantBytes(t, got, want)
	})

	t.Run("header", func(t *testing.T) {
		path := "/xxx/yyy"
		req, _ := http.NewRequest("GET", url+path, nil)
		req.Header.Set("Authorization", "Bearer access-token")
		got := toPathsWithValidation(t, req, nil)
		want := []byte(`
{
  "/xxx/yyy": {
    "get": {
      "parameters": [
        {
          "in": "header",
          "name": "Authorization",
          "examples": [
            "Bearer access-token"
          ]
        }
      ]
    }
  }
}
`)
		difftest.AssertGotAndWantBytes(t, got, want)
	})
}

func toPathsWithValidation(t *testing.T, req *http.Request, body io.Reader) Paths {
	t.Helper()
	paths, err := toPaths(req, body)
	if err != nil {
		t.Fatalf("toPath failed, %+v", err)
	}
	return paths
}
