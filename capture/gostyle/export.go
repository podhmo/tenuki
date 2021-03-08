package gostyle

import (
	"fmt"
	"net/http"

	"github.com/podhmo/tenuki/capture/style"
)

func ExtractRequestInfo(req *http.Request) (style.Info, error) {
	raw := Info{}
	info, err := parseRequest(req)
	if err != nil {
		return raw, fmt.Errorf("parse request %w", err)
	}
	raw["Request"] = info
	return raw, nil
}

func ExtractResponseInfo(resp *http.Response, info style.Info) (style.Info, error) {
	raw, ok := info.(Info)
	if !ok {
		return info, fmt.Errorf("unexpected info %T", info)
	}
	info, err := parseResponse(resp)
	if err != nil {
		return raw, fmt.Errorf("parse response %w", err)
	}
	raw["Response"] = info
	return raw, nil
}
