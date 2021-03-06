package tenuki_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/podhmo/tenuki"
)

func TestHandlerRoundTripper(t *testing.T) {
	transport := tenuki.HandlerTripperFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world"+r.URL.Query().Get("suffix"))
	})

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatalf("! %+v", err)
	}

	q := req.URL.Query()
	q.Add("suffix", " !!")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{
		Timeout:   100 * time.Millisecond,
		Transport: transport,
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("!! %+v", err)
	}

	want := "hello world !!"

	var b bytes.Buffer
	if _, err := b.ReadFrom(res.Body); err != nil {
		t.Fatalf("!!! %+v", err)
	}
	defer res.Body.Close()

	got := b.String()
	if want != got {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}

func TestErrorTransport(t *testing.T) {
	thisErr := fmt.Errorf("THIS!!")
	transport := tenuki.NewErrorTransport(t, func() error { return thisErr })
	req, _ := http.NewRequest("GET", "", nil)

	client := &http.Client{
		Timeout:   100 * time.Millisecond,
		Transport: transport,
	}

	_, err := client.Do(req)
	if !errors.Is(err, thisErr) {
		t.Errorf("unexpected error is returned, %[1]T, %+[1]v", err)
	}
}
