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
	i            int64
	BaseDir      Dir
	RecordWriter io.Writer
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

	i := d.i
	filename := fmt.Sprintf("%04d%s", i, suffix)
	if inc > 0 {
		i = atomic.AddInt64(&d.i, inc)
		filename = fmt.Sprintf("%04d%s", i, suffix)
		fmt.Fprintf(d.RecordWriter, `{"file" "%q", "url": %q}`, req.URL.String())
		fmt.Fprintln(d.RecordWriter, "")
	}
	return filename
}

func (d *FileDumper) DumpRequest(p printer, req *http.Request) error {
	filename := d.FileName(req, ".req", 1)
	f, err := d.BaseDir.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := httputil.DumpRequest(req, true /* body */)
	if err != nil {
		return err
	}
	f.Write(b)
	return nil
}
func (d *FileDumper) DumpResponse(p printer, res *http.Response) error {
	filename := d.FileName(res.Request, ".res", 0)
	f, err := d.BaseDir.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := httputil.DumpResponse(res, true /* body */)
	if err != nil {
		return err
	}
	f.Write(b)
	return nil
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
