package openapistyle

import (
	"fmt"
	"io"
	"net/http"

	"github.com/podhmo/tenuki/capture/style"
)

type Info struct {
	Method      string `json:"method"`
	URL         string `json:"url"`
	HTTPVersion string `json:"httpVersion,omitempty"`
	HeaderSize  int64  `json:"headerSize,omitempty" default:"-1"`
	BodySize    int64  `json:"bodySize,omitempty" default:"-1"`

	ContentType string `json:"contentType"`
	Paths       Paths  `json:"paths"`
}

func (info *Info) Merge(res style.Info) style.Info {
	// TODO
	return info
}

// TODO
func (info *Info) HandleError(open func() (io.WriteCloser, error), err error) {
}

func ExtractRequestInfo(req *http.Request) (style.Info, error) {
	info := Info{}
	paths, err := toPaths(req, &info)
	if err != nil {
		return nil, fmt.Errorf("extract paths, %w", err)
	}
	info.Paths = paths
	return &info, nil
}

func ExtractResponseInfo(resp *http.Response) (style.Info, error) {
	// TODO:
	return &Info{}, nil
}

// Return value if nonempty, def otherwise.
func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}
