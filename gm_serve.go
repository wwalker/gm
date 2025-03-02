package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kpym/gm/internal/browser"
)

// availablePort provides the first available port after 8080
// or 8180 if no available ports are present.
func availablePort() (port string) {
	for i := 8080; i < 8181; i++ {
		port = strconv.Itoa(i)
		if ln, err := net.Listen("tcp", "localhost:"+port); err == nil {
			ln.Close()
			break
		}
	}
	return port
}

// serveFiles serve the local folder `serveDir`.
// If an .md (or corresponding .html) file is requested it is compiled and send as html.
func serveFiles() {
	var lastMethodPath string

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join(serveDir, r.URL.Path)
		newMethodPath := fmt.Sprintf("\n%s '%s':", r.Method, r.URL.Path)
		if newMethodPath != lastMethodPath {
			lastMethodPath = newMethodPath
			info(newMethodPath)
		}
		if strings.HasSuffix(filename, ".html") {
			filename = filename[0:len(filename)-5] + ".md"
		}
		if strings.HasSuffix(filename, "md") {
			if r.Method == "HEAD" {
				info(".")
				if fstat, err := os.Stat(filename); err == nil {
					w.Header().Set("Last-Modified", fstat.ModTime().UTC().Format(http.TimeFormat))
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte{})
				}
				return
			}
			if content, err := ioutil.ReadFile(filename); err == nil {
				if html, err := compile(content); err == nil {
					info(" serve converted .md file.")
					w.Write(html)
					return
				}
			}
		}
		if r.URL.Path == "/favicon.ico" {
			info(" serve internal png.")
			w.Header().Set("Cache-Control", "max-age=86400") // 86400 s = 1 day
			w.Header().Set("Expires", time.Now().Add(24*time.Hour).UTC().Format(http.TimeFormat))
			w.Write(favIcon)
			return
		}
		info(" serve raw file.")
		w.Header().Set("Cache-Control", "no-store")
		http.FileServer(http.Dir(serveDir)).ServeHTTP(w, r)
	})

	port := availablePort()
	info("start serving '%s' folder to localhost:%s.\n", serveDir, port)
	url := "http://localhost:" + port + "/" + serveFile
	err := browser.Open(url)
	try(err, "Can't open the web browser, but you can visit now:", url)
	check(http.ListenAndServe("localhost:"+port, nil))
}
