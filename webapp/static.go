package main

import (
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var etag = make(map[string]string)

var expiration = map[string]time.Duration{
	"/fonts/": 365 * 24 * time.Hour,
	"/png/":   1 * time.Hour,
	"/img/":   1 * time.Hour,
	"/js/":    10 * time.Minute,
	"/css/":   10 * time.Minute,
}

var rootdir = regexp.MustCompile(`^/([^/]*)/`)

func ServeFile(w http.ResponseWriter, req *http.Request, filename string) {
	if tag, ok := etag[filename]; ok {
		if req.Header.Get("If-None-Match") == tag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	file, err1 := os.Open(filename)
	if err1 != nil {
		http.NotFound(w, req)
		return
	}
	defer file.Close()

	info, err2 := file.Stat()
	if err2 != nil {
		http.NotFound(w, req)
		return
	}

	// ETag
	hash := sha1.New()
	io.Copy(hash, file)
	etag[filename] = fmt.Sprintf("%x", hash.Sum(nil))
	file.Seek(0, os.SEEK_SET)

	if req.Header.Get("If-None-Match") == etag[filename] {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// expire depending on directory
	expDur := 1 * time.Hour
	dir := rootdir.FindString(req.URL.Path)
	log.Printf("DIR: %s", dir)
	if dur, ok := expiration[dir]; ok {
		expDur = dur
	}
	expDate := time.Now().Add(expDur)

	w.Header().Set("ETag", etag[filename])
	w.Header().Set("Expires", expDate.UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))

	if enc := req.Header.Get("Accept-Encoding"); strings.Contains(enc, "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		io.Copy(gz, file)
	} else {
		w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
		io.Copy(w, file)
	}
}

func hStaticFiles(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s\n", req.Method, req.URL)
	filename := filepath.Join(srvdir, "static"+req.URL.Path)
	ServeFile(w, req, filename)
}
