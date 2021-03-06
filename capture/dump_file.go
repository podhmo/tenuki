package capture

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"sync/atomic"
)

type FileDumper struct {
	BaseDir Dir

	Counter      *int64
	DisableCount bool

	RecordWriter io.Writer
	Prefix       string
}

func (d *FileDumper) FileName(req *http.Request, suffix string, inc int64) string {
	if d.RecordWriter == nil {
		f, err := d.BaseDir.Open("records.txt")
		// xxx: does not Close()

		d.RecordWriter = f
		if err != nil {
			if err != nil {
				log.Printf("create records.txt failured: %+v", err)
			}
			d.RecordWriter = ioutil.Discard
		}
	}

	prefix := d.Prefix
	if d.Counter == nil {
		n := int64(0)
		d.Counter = &n
	}
	if !d.DisableCount {
		i := atomic.AddInt64(d.Counter, inc)
		prefix = fmt.Sprintf("%04d%s", i, prefix)
	}
	filename := prefix + suffix

	url := "/"
	if req != nil {
		url = req.URL.String()
	}
	fmt.Fprintf(d.RecordWriter, "{\"file\": %q, \"url\": %q}\r\n", filename, url)
	return filename
}

func (d *FileDumper) DumpRequest(p printer, req *http.Request) (State, error) {
	filename := d.FileName(req, ".req", 1)
	state := fileState{request: req, FileName: filename}
	f, err := d.BaseDir.Open(filename)
	if err != nil {
		return state, err
	}
	defer f.Close()

	b, err := httputil.DumpRequest(req, true /* body */)
	if err != nil {
		return state, err
	}
	f.Write(b)
	return state, nil
}

func (d *FileDumper) DumpError(p printer, state State, err error) error {
	req := state.Request()
	filename := d.FileName(req, ".error", 0)
	f, _ := d.BaseDir.Open(filename)
	d.dumpHeader(f, req)
	fmt.Fprintf(f, "%+v\n", err)
	return err
}

func (d *FileDumper) DumpResponse(p printer, state State, res *http.Response) error {
	req := res.Request
	filename := d.FileName(req, ".res", 0)
	f, err := d.BaseDir.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if req != nil {
		d.dumpHeader(f, req)
	}
	b, err := httputil.DumpResponse(res, true /* body */)
	if err != nil {
		return err
	}
	f.Write(b)
	return nil
}

func (d *FileDumper) dumpHeader(w io.Writer, req *http.Request) {
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

type fileState struct {
	request  *http.Request
	FileName string
}

func (s fileState) Request() *http.Request {
	return s.request
}

var _ Dumper = &FileDumper{}
