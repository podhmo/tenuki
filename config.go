package tenuki

import (
	"log"
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
	writeFileBaseDir string
	layout           *capture.Layout
}

func DefaultConfig(options ...func(*Config)) *Config {
	c := &Config{
		layout:           DefaultLayout,
		captureEnabled:   CaptureEnabledDefault,
		writeFileBaseDir: CaptureWriteFileBaseDir,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c *Config) NewCaptureTransport(prefix string) *capture.CapturedTransport {
	ct := &capture.CapturedTransport{
		Printer: log.New(os.Stderr, "tenuki", 0),
		Dumper:  &capture.ConsoleDumper{Layout: c.layout},
	}
	if c.writeFileBaseDir != "" {
		ct.Dumper = &capture.FileDumper{
			FileManager: getFileManagerWithDefault(c.writeFileBaseDir),
			Layout:      c.layout,
			Prefix:      prefix,
		}
	}
	return ct
}

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
