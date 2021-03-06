package gostyle

import (
	"io"
	"net/http"
)

func ExtractRequestInfo(req *http.Request, body io.Reader) (interface{ Info() interface{} }, error) {
	return parseRequest(req, body)
}

func ExtractResponseInfo(resp *http.Response, body io.Reader) (interface{ Info() interface{} }, error) {
	return parseResponse(resp, body)
}
