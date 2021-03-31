package web

import (
	"embed"
	_ "embed"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed script/ipyshell-console.py
var ipyshellconsole []byte

//go:embed dist/static
var static embed.FS

//go:embed dist/index.html
var index []byte

//go:embed dist/widget.html
var widget []byte

//go: embed dist/favicon.ico
var favicon []byte

type embedFS struct {
	embed.FS
	path string
}
type embedFile struct {
	io.Seeker
	fs.File
}

func (*embedFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}

func (fs *embedFS) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	name = strings.Split(name, "?")[0]
	fullName := filepath.Join(fs.path, filepath.FromSlash(path.Clean("/"+name)))
	file, err := fs.FS.Open(fullName)
	ef := &embedFile{
		File: file,
	}
	return ef, err
}

// Static Static
func Static() http.FileSystem {
	if err := recover(); err != nil {
		//fmt.Println(err)
	}
	return &embedFS{
		static,
		"dist/static",
	}
}

// Index Index
func Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", index)
	}
}

// Widget Widget
func Widget() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", widget)
	}
}

// Favicon Favicon
func Favicon() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "image/ico", favicon)
	}
}

func TempScript() string {
	fname := filepath.Join(os.TempDir(), "ipyshell-console.py")
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		os.WriteFile(fname, ipyshellconsole, 0755)
	}
	return fname
}
