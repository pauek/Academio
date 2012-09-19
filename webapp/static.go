package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var etag = make(map[string]string)

func StaticFiles(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s\n", req.Method, req.URL)
	filename := filepath.Join(srvdir, "static"+req.URL.Path)

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
	w.Header().Set("ETag", etag[filename])
	oneHourFromNow := time.Now().Add(1 * time.Hour)
	w.Header().Set("Expires", oneHourFromNow.UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(w, file)
}
