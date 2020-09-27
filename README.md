# tenuki

test helper for web api testing

## features

- api request helper
- capture transport

### api request helper

```go
func Test(t *testing.T){
    handler := func(w http.ResponseWriter, r *http.Request){
		tenuki.Render(w, r).JSON(200, map[string]string{"message": "hello world"})
    }
	ts := httptest.NewServer(http.HandlerFunc(handler))

    f := tenuki.New(t)
	req := f.NewRequest("GET", "http://localhost:8080/hello", nil)
	res := f.Do(req)

	if want, got := 200, res.StatusCode; want != got {
		t.Errorf("status code\nwant\n\t%d\nbut\n\t%d", want, got)
	}

	var got map[string]string{}
	f.Extract().JSON(res, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
```

Requesting via `tenkuki.Facade.Do()`, capture transport is automatically activated.
(If stopping this feature, please go test with `NOCAPTURE=1` envvar)

### capture transport

output example.

```
=== RUN   TestCapture/request_2
    request:
        POST / HTTP/1.1
        Host: 127.0.0.1:59424
        
        {"me": "foo"}
    response:
        HTTP/1.1 200 OK
        Content-Length: 21
        Content-Type: text/plain; charset=utf-8
        Date: Sun, 27 Sep 2020 14:14:55 GMT
        
        {"message": "hello"}
```

code example.

```go
func TestCapture(t *testing.T) {
	transport := &tenuki.CapturedTransport{}

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"message": "hello"}`)
		},
	))
	defer ts.Close()

	client := &http.Client{Transport: transport}

    // GET request
	t.Run("request 1", func(t *testing.T) {
		defer transport.Capture(t)()

		client.Get(ts.URL)
	})

    // POST request
	t.Run("request 2", func(t *testing.T) {
		defer transport.Capture(t)()

		req, _ := http.NewRequest("POST", ts.URL, strings.NewReader(`{"me": "foo"}`))
		client.Do(req)
	})
}
```
