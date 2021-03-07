package gostyle

import (
	"io"
	"net/http"

	"github.com/podhmo/tenuki/capture/style"
)

func ExtractRequestInfo(req *http.Request, body io.Reader) (style.Info, error) {
	return parseRequest(req, body)
}

func ExtractResponseInfo(resp *http.Response, body io.Reader) (style.Info, error) {
	return parseResponse(resp, body)
}
