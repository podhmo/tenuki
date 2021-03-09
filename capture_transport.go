package tenuki

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

type CapturedTransport struct {
	T         *testing.T
	Transport http.RoundTripper
}

func (ct *CapturedTransport) Printf(fmt string, args ...interface{}) {
	ct.T.Helper()
	ct.T.Logf("\x1b[5G\x1b[0K"+fmt, args...)
}

func (ct *CapturedTransport) GetPrefix() string {
	return ct.T.Name()
}

func (ct *CapturedTransport) Capture(t *testing.T) func() {
	ct.T = t
	return func() {
		ct.T = nil
	}
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if ct.T == nil {
		fmt.Fprintln(os.Stderr, "!! CapturedTransport.T is not found !!")
		fmt.Fprintln(os.Stderr, "please use `defer transport.Capture(t)()`")
	}
	return ct.Transport.RoundTrip(req)
}

func NewCaptureTransport(t *testing.T, transport http.RoundTripper, options ...func(*Config)) *CapturedTransport {
	c := DefaultConfig(append([]func(*Config){func(c *Config) {
		c.disableCount = true
	}}, options...)...)
	ct := &CapturedTransport{T: t}
	ct.Transport = c.NewCaptureTransport(transport, ct.GetPrefix)
	return ct
}
