package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPlainOld(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(map[string]string{"message": "hello world"}); err != nil {
			// handling err
			panic(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))

	req, err := http.NewRequest("GET", ts.URL+"/hello", nil)
	if err != nil {
		t.Fatalf("req: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("res: %v", err)
	}

	if want, got := 200, res.StatusCode; want != got {
		t.Errorf("status code\nwant\n\t%d\nbut\n\t%d", want, got)
	}

	want := map[string]string{"message": "hello world"}
	var got map[string]string
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	defer res.Body.Close()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
