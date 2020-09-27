package reqtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
)

type RenderFacade struct {
	w  http.ResponseWriter
	r  *http.Request
	ok int
	ng int
}

func Render(w http.ResponseWriter, r *http.Request) *RenderFacade {
	return &RenderFacade{
		w:  w,
		r:  r,
		ok: http.StatusOK,
		ng: http.StatusInternalServerError,
	}
}
func (f *RenderFacade) SetOKStatus(code int) *RenderFacade {
	f.ok = code
	return f
}
func (f *RenderFacade) SetNGStatus(code int) *RenderFacade {
	f.ng = code
	return f
}
func (f *RenderFacade) JSON(v interface{}) {
	w := f.w

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(v); err != nil {
		http.Error(w, err.Error(), f.ng)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if f.ok > 0 {
		w.WriteHeader(f.ok)
	}
	w.Write(buf.Bytes())
}
func (f *RenderFacade) JSONArray(v interface{}) {
	// Force to return empty JSON array [] instead of null in case of zero slice.
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice && val.IsNil() {
		v = reflect.MakeSlice(val.Type(), 0, 0).Interface()
	}

	f.JSON(v)
}
