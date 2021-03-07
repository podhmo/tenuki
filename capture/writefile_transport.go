package capture

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
)

type FileManager struct {
	BaseDir Dir

	Counter      *int64
	DisableCount bool

	RecordWriter io.Writer
}

func (m *FileManager) FileName(req *http.Request, name string, suffix string, inc int64) string {
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
	err := wt.DumpRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, wt.DumpError(req, err)
	}
	if err := wt.DumpResponse(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (wt *WriteFileTransport) DumpRequest(req *http.Request) error {
	filename := wt.FileName(req, wt.GetPrefix(), ".req", 1)
	f, err := wt.BaseDir.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	layout := wt.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	b, err := layout.Request.Extract(req)
	if err != nil {
		return err
	}

	f.Write(b)
	return nil
}

func (wt *WriteFileTransport) DumpError(req *http.Request, err error) error {
	filename := wt.FileName(req, wt.GetPrefix(), ".error", 0)
	f, _ := wt.BaseDir.Open(filename)
	wt.dumpHeader(f, req)
	fmt.Fprintf(f, "%+v\n", err)
	return err
}

func (wt *WriteFileTransport) DumpResponse(req *http.Request, res *http.Response) error {
	filename := wt.FileName(req, wt.GetPrefix(), ".res", 0)
	f, err := wt.BaseDir.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if req != nil {
		wt.dumpHeader(f, req)
	}
	layout := wt.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	b, err := layout.Response.Extract(res)
	f.Write(b)
	return nil
}

func (wt *WriteFileTransport) dumpHeader(w io.Writer, req *http.Request) {
	reqURI := req.RequestURI
	if reqURI == "" {
		reqURI = req.URL.RequestURI()
	}
	method := req.Method
	if method == "" {
		method = "GET"
	}
	fmt.Fprintf(w, "%s %s HTTP/%d.%d\r\n", method,
		reqURI, req.ProtoMajor, req.ProtoMinor)
}
