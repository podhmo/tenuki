package gostyle

import (
	"net/http"

	"github.com/podhmo/tenuki/capture/style"
)

func ExtractRequestInfo(req *http.Request) (style.Info, error) {
	info, err := parseRequest(req)
	return info, err
}

func ExtractResponseInfo(resp *http.Response) (style.Info, error) {
	info, err := parseResponse(resp)
	return info, err
}
