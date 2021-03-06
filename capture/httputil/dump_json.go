package httputil

import (
	"net/http"
)

type RequestInfo struct {
	URI    string
	Method string
	Proto  string
	Body   []byte
}

func DumpRequestJSON(req *http.Request, body bool) (RequestInfo, error) {
	var err error
	var info RequestInfo

	save := req.Body
	{
		if !body || req.Body == nil {
			req.Body = nil
		} else {
			save, req.Body, err = drainBody(req.Body)
			if err != nil {
				return info, err
			}
		}
	}

	_ = save
	return info, nil
}

type ResponseInfo struct {
}

func DumpResponseJSON(resp *http.Response, body bool) (ResponseInfo, error) {
	info := ResponseInfo{}

	var err error
	save := resp.Body
	savecl := resp.ContentLength

	{
		if !body {
			// For content length of zero. Make sure the body is an empty
			// reader, instead of returning error through failureToReadBody{}.
			if resp.ContentLength == 0 {
				resp.Body = emptyBody
			} else {
				resp.Body = failureToReadBody{}
			}
		} else if resp.Body == nil {
			resp.Body = emptyBody
		} else {
			save, resp.Body, err = drainBody(resp.Body)
			if err != nil {
				return info, err
			}
		}

		if err == errNoBody {
			err = nil
		}
		resp.Body = save
		resp.ContentLength = savecl
	}

	if err != nil {
		return info, err
	}
	return info, nil
}
