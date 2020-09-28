package capture

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

type FileDumper struct {
	i       int64
	BaseDir Dir
}

func (d *FileDumper) FileName(req *http.Request, suffix string, inc int64) string {
	i := atomic.AddInt64(&d.i, inc)
	method := req.Method
	return fmt.Sprintf("%04d%s@%s", i, method, strings.Replace(req.URL.String(), "/", "@", -1)+suffix)
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

	log.Println("\ttrace to", filename)
	return os.Create(filepath.Join(dir, filename))
}
