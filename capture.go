package tenuki

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/podhmo/tenuki/capture"
)

var (
	fileManagerMap = map[string]*capture.FileManager{}
	mu             sync.Mutex
)

func getFileManagerWithDefault(basedir string) *capture.FileManager {
	mu.Lock()
	defer mu.Unlock()

	m, ok := fileManagerMap[basedir]
	if ok {
		return m
	}
	var c int64
	m = &capture.FileManager{
		BaseDir:      capture.Dir(basedir),
		DisableCount: !CaptureCountEnabledDefault,
		Counter:      &c,
	}
	fileManagerMap[basedir] = m
	return m
}

func NewCaptureTransportWithDefault(t *testing.T, basedir string) *CapturedTransport {
	ct := &CapturedTransport{T: t}
	if basedir != "" {
		ct.Dumper = &capture.FileDumper{
			FileManager: getFileManagerWithDefault(basedir),
			Prefix:      t.Name(),
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
