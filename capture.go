package tenuki

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/podhmo/tenuki/capture"
)

func NewCaptureTransportWithDefault(t *testing.T, basedir string) *CapturedTransport {
	ct := &CapturedTransport{T: t}
	if basedir != "" {
		ct.Dumper = &capture.FileDumper{
			BaseDir:      capture.Dir(basedir),
			DisableCount: !CaptureCountEnabledDefault,
			Counter:      &globalFileDumpCounter,
			Prefix:       t.Name(),
		}
	}
	return ct
}

type CapturedTransport struct {
	capture.CapturedTransport
	T *testing.T
}

func (ct *CapturedTransport) Capture(t *testing.T) func() {
	ct.T = t
	return func() {
		ct.T = nil
	}
}

func (ct *CapturedTransport) Printf(fmt string, args ...interface{}) {
	ct.T.Helper()
	ct.T.Logf("\x1b[5G\x1b[0K"+fmt, args...)
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if ct.T == nil {
		fmt.Fprintln(os.Stderr, "!! CapturedTransport.T is not found !!")
		fmt.Fprintln(os.Stderr, "please use `defer transport.Capture(t)()`")
	}
	if ct.CapturedTransport.Printer == nil {
		ct.CapturedTransport.Printer = ct
	}
	return ct.CapturedTransport.RoundTrip(req)
}
