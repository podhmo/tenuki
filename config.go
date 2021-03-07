package tenuki

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/podhmo/tenuki/capture"
)

var (
	CaptureEnabledDefault   bool   = true
	CaptureWriteFileBaseDir string = ""

	CaptureCountEnabledDefault bool = false
	DefaultLayout              *capture.Layout
)

func init() {
	if ok, _ := strconv.ParseBool(os.Getenv("NOCAPTURE")); ok {
		log.Println("CAPTURE_DISABLED is true, so deactivate tenuki.capture function")
		CaptureEnabledDefault = false
	}
	if ok, _ := strconv.ParseBool(os.Getenv("CAPTURE_DISABLED")); ok {
		log.Println("CAPTURE_DISABLED is true, so deactivate tenuki.capture function")
		CaptureEnabledDefault = false
	}
	if filename := os.Getenv("CAPTURE_WRITEFILE"); filename != "" {
		log.Println("CAPTURE_WRITEFILE is set, so activate the function writing capture output to files")
		CaptureWriteFileBaseDir = filename
	}

	if layout := os.Getenv("CAPTURE_LAYOUT"); layout != "" {
		switch strings.ToLower(layout) {
		case "text":
			DefaultLayout = capture.TextLayout
		case "json":
			DefaultLayout = capture.JSONLayout
		case "openapi":
			DefaultLayout = capture.OpenAPILayout
		default:
			log.Printf("layout=%q is not found, use text layout. (availables: text, json, openapi)", layout)
			DefaultLayout = capture.TextLayout
		}
	}
}

type Config struct {
	captureEnabled   bool
	disableCount     bool
	writeFileBaseDir string

	layout  *capture.Layout
	printer printer
}
type printer interface {
	Printf(fmt string, args ...interface{})
}

func DefaultConfig(options ...func(*Config)) *Config {
	c := &Config{
		layout:           DefaultLayout,
		captureEnabled:   CaptureEnabledDefault,
		writeFileBaseDir: CaptureWriteFileBaseDir,
		printer:          log.New(os.Stderr, "tenuki", 0),
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func WithoutCapture() func(*Config) {
	return func(c *Config) {
		c.captureEnabled = false
	}
}
func WithWriteFile(basedir string) func(*Config) {
	return func(c *Config) {
		c.writeFileBaseDir = basedir
	}
}
func WithLayout(layout *capture.Layout) func(*Config) {
	return func(c *Config) {
		c.layout = layout
	}
}
func WithPrinter(printer printer) func(*Config) {
	return func(c *Config) {
		c.printer = printer
	}
}

func (c *Config) NewCaptureTransport(transport http.RoundTripper, getPrefix func() string) http.RoundTripper {
	if transport == nil {
		transport = http.DefaultTransport
	}
	switch c.writeFileBaseDir {
	case "":
		return &capture.ConsoleTransport{
			Transport: transport,
			Printer:   c.printer,
			Layout:    c.layout,
		}
	default:
		return &capture.WriteFileTransport{
			Transport:   transport,
			FileManager: getFileManagerWithDefault(c.writeFileBaseDir, c.disableCount),
			Layout:      c.layout,
			GetPrefix:   getPrefix,
		}
	}
}

var (
	fileManagerMap = map[string]*capture.FileManager{}
	mu             sync.Mutex
)

func getFileManagerWithDefault(basedir string, disableCount bool) *capture.FileManager {
	mu.Lock()
	defer mu.Unlock()

	m, ok := fileManagerMap[basedir]
	if ok {
		return m
	}
	var c int64
	m = &capture.FileManager{
		BaseDir:      capture.Dir(basedir),
		DisableCount: disableCount,
		Counter:      &c,
	}
	fileManagerMap[basedir] = m
	return m
}
