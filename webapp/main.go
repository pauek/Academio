package main

import (
	"Academio/content"
	"Academio/webapp/data"
	"flag"
	"fmt"
	F "fragments"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var srvdir = filepath.Join(os.Getenv("ACADEMIO_ROOT"), "webapp")
var cache = F.NewCache()

func hFragList(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/_frag/"):]
	list := cache.List("item " + id)
	w.Header().Set("Content-Type", "text/html")
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

func hLogin(w http.ResponseWriter, req *http.Request) {
	session := data.GetOrCreateSession(req)
	switch req.Method {
	case "GET":
		url, err := url.Parse(req.Header.Get("Referer"))
		if err == nil && url.Path != "/login" {
			session.Referer = url.Path
			log.Printf("Referer = '%s'", url.Path)
		}
		FragmentDispatch(w, req, session, "login", "Login")
	case "POST":
		hLoginProcessForm(w, req, session)
	default:
		http.Error(w, "Wrong method", http.StatusBadRequest)
	}
}

func hLoginProcessForm(w http.ResponseWriter, req *http.Request, session *data.Session) {
	login := req.FormValue("login")
	password := req.FormValue("password")
	if user := data.AuthenticateUser(login, password); user != nil {
		session.SetUser(user)
		url := session.Referer
		if url == "" {
			url = "/"
		}
		http.Redirect(w, req, url, http.StatusSeeOther)
		return
	}
	session.Message = "Incorrect Login"
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func hLogout(w http.ResponseWriter, req *http.Request) {
	session := data.GetSession(req)
	if session != nil {
		session.User = nil
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
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
	http.HandleFunc("/login", hLogin)
	http.HandleFunc("/logout", hLogout)
	http.HandleFunc("/", GzippedFunc(fragmentPage))

	if *ssl {
		go listenSSL()
		redirectToSSL()
	} else {
		listen()
	}
}
