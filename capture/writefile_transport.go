package capture

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/podhmo/tenuki/capture/style"
)

type WriteFileTransport struct {
	Transport http.RoundTripper
	*FileManager
	Layout    *Layout
	GetPrefix func() string // xxx: use t.Name()
}

func (wt *WriteFileTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := wt.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	s, err := wt.DumpRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, wt.DumpError(req, s, err)
	}
	if err := wt.DumpResponse(res, req, s); err != nil {
		return nil, err
	}
	return res, nil
}

func (wt *WriteFileTransport) DumpRequest(req *http.Request) (style.State, error) {
	layout := wt.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	s, err := layout.Request.Extract(req)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (wt *WriteFileTransport) DumpError(req *http.Request, s style.State, err error) error {
	s.Info().HandleError(func() (io.WriteCloser, error) {
		filename := wt.FileName(req, wt.GetPrefix(), ".error", 0)
		return wt.BaseDir.Open(filename)
	}, err)
	return err
}

func (wt *WriteFileTransport) DumpResponse(res *http.Response, req *http.Request, s style.State) error {
	filename := wt.FileName(req, wt.GetPrefix(), ".res", 0)
	f, err := wt.BaseDir.Open(filename)
	if err != nil {
		return fmt.Errorf("in response, open: %w", err)
	}
	defer f.Close()

	// if req != nil {
	// 	wt.dumpHeader(f, req)
	// }
	layout := wt.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	s2, err := layout.Response.Extract(res, s)
	if err := s2.Emit(f); err != nil {
		return err
	}
	return nil
}

type FileManager struct {
	BaseDir Dir

	Counter      *int64
	DisableCount bool

	RecordWriter io.Writer
}

func (m *FileManager) FileName(req *http.Request, name string, suffix string, inc int64) string {
	if strings.Contains(name, "/") {
		name = strings.ReplaceAll(name, "/", "__")
	}

	if m.RecordWriter == nil {
		f, err := m.BaseDir.Open("RECORDS.txt")
		// xxx: does not Close()

		m.RecordWriter = f
		if err != nil {
			if err != nil {
				log.Printf("create RECORDS.txt failured: %+v", err)
			}
			m.RecordWriter = ioutil.Discard
		}
	}

	if m.Counter == nil {
		n := int64(0)
		m.Counter = &n
	}
	prefix := ""
	if !m.DisableCount {
		i := atomic.AddInt64(m.Counter, inc)
		prefix = fmt.Sprintf("%04d", i)
	}
	method := "GET"
	if req != nil {
		method = req.Method
	}

	filename := fmt.Sprintf("%s%s@%s%s", prefix, name, method, suffix)
	url := "/"
	if req != nil {
		url = req.URL.String()
	}
	fmt.Fprintf(m.RecordWriter, "{\"file\": %q, \"url\": %q}\r\n", filename, url)
	return filename
}

type Dir string

func (d Dir) Open(filename string) (io.WriteCloser, error) {
	dir := string(d)
	if dir != "" {
		if err := os.MkdirAll(dir, 0744); err != nil {
			return nil, err
		}
	}

	fullname := filepath.Join(dir, filename)
	log.Println("\ttrace to", fullname)
	return os.Create(fullname)
}
