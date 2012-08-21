package main

import (
	"Academio/content"
	"flag"
	"fmt"
	F "fragments"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var srvdir = filepath.Join(os.Getenv("ACADEMIO_ROOT"), "webapp")
var cache = F.NewCache()

func hFragList(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/_frag/"):]
	list := cache.List("item " + id)
	tmpl.Lookup("fraglist").Execute(w, list)
}

func fmtTime(t time.Time) string {
	return t.UTC().Format(http.TimeFormat)
}

func hPhotos(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/png/"):]

	item := content.Get(id)
	if item == nil {
		http.NotFound(w, req)
		return
	}
	course, ok := item.(*content.Course)
	if !ok {
		http.NotFound(w, req)
		return
	}

	image, err := os.Open(course.Photo)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	defer image.Close()

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Expires", fmtTime(time.Now().Add(1 * time.Hour)))
	if stat, err := image.Stat(); err == nil {
		w.Header().Set("Last-Modified", fmtTime(stat.ModTime()))
	}
	if req.Method == "HEAD" {
		return
	}
	io.Copy(w, image)
}

var port = flag.Int("port", 8080, "Network port")
var ssl = flag.Bool("ssl", false, "Use SSL?")

func serveFiles(prefix string) {
	fs := http.FileServer(http.Dir(srvdir + prefix))
	http.Handle(prefix, http.StripPrefix(prefix, GzippedNoExpire(fs)))
}

func listen() {
	p := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(p, nil)
	if err != nil {
		log.Fatalf("Cannot Listen: %s", err)
	}
}

func listenSSL() {
	root := os.Getenv("ACADEMIO_ROOT")
	certfile := filepath.Join(root, "webapp/certs/cert.pem")
	keyfile  := filepath.Join(root, "webapp/certs/academio.key")
		p := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServeTLS(p, certfile, keyfile, nil)
	if err != nil {
		log.Fatalf("Cannot ListenTLS: %s", err)
	}
}

func redirectToSSL() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		url := "https:/" + "/academ.io" + req.URL.String()
		http.Redirect(w, req, url, http.StatusMovedPermanently)
	})
	srv := http.Server{
		Addr:    ":http",
		Handler: mux,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot Listen (http -> https redirect): %s", err)
	}
}

func main() {
	flag.Parse()

	content.WatchForChanges(func(id string) {
		if id == "" {
			cache.Touch("/courses")
		} else {
			cache.Touch("/content/" + id)
		}
	})

	// handlers
	serveFiles("/js/lib/")
	serveFiles("/js/")
	serveFiles("/css/")
	serveFiles("/img/")

	http.HandleFunc("/_frag/", GzippedFunc(hFragList))
	http.HandleFunc("/png/", hPhotos)
	http.HandleFunc("/", GzippedFunc(fragmentPage))

	if *ssl {
		go listenSSL()
		redirectToSSL()
	} else {
		listen()
	}
}
