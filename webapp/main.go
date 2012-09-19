package main

import (
	"Academio/content"
	"Academio/webapp/data"
	"bytes"
	"flag"
	"fmt"
	F "fragments"
	"html/template"
	"io"
	"log"
	"mime"
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
	var zero time.Time
	list := cache.Diff("item " + id, zero)
	w.Header().Set("Content-Type", "text/html")
	tmpl.Lookup("fraglist").Execute(w, list)
}

func fmtTime(t time.Time) string {
	return t.UTC().Format(http.TimeFormat)
}

func NotFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html")
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "notfound", req.URL.Path); err != nil {
		goto InternalError
	}
	tmpl.ExecuteTemplate(w, "layout", layoutInfo{
		Title:   "No encontrada - Academio",
		Message: "",
		Navbar:  template.HTML(cache.RenderToString("navbar")),
		Body:    template.HTML(body.String()),
	})
	return
	
InternalError:
	code := http.StatusInternalServerError
	http.Error(w, "Template 'layout' not found", code)
}


func hFavicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, filepath.Join(srvdir, "static/img/favicon.ico"))
}

func REMOVE_hFavicon(w http.ResponseWriter, req *http.Request) {
	rest := req.URL.Path[len("/favicon.ico"):]
	if len(rest) > 0 {
		log.Printf("%s (NOT FOUND)", req.URL.Path)
		NotFound(w, req)
		return
	}
	path := filepath.Join(srvdir, "static/img/favicon.ico")
	icon, err := os.Open(path)
	if err != nil {
		NotFound(w, req)
		return
	}
	defer icon.Close()
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Expires", fmtTime(time.Now().Add(24 * time.Hour)))
	io.Copy(w, icon)
}

func hPhotos(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/png/"):]
	item := content.Get(id)
	if item == nil {
		NotFound(w, req)
		return
	}
	course, ok := item.(*content.Course)
	if !ok {
		NotFound(w, req)
		return
	}
	http.ServeFile(w, req, course.Photo)
}

func hLogin(w http.ResponseWriter, req *http.Request) {
	session := data.GetOrCreateSession(req)
	log.Printf("%s [%s]", req.URL, session.Id)
	session.PutCookie(w)
	switch req.Method {
	case "GET":
		url, err := url.Parse(req.Header.Get("Referer"))
		if err == nil && url.Path != "/login" {
			session.Referer = url.Path
		}
		SendPage(w, req, session, "login", "Login")
		session.Message = ""
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
	session.Message = "Login incorrecto"
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func hLogout(w http.ResponseWriter, req *http.Request) {
	session := data.GetSession(req)
	log.Printf("%s [%s]", req.URL, session.Id)
	if session != nil {
		session.User = nil
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func hRegister(w http.ResponseWriter, req *http.Request) {
	session := data.GetSession(req)
	log.Printf("%s [%s]", req.URL, session.Id)
	if req.URL.Path != "/register" {
		NotFound(w, req)
		return
	}
	NotFound(w, req)
}

func hAbout(w http.ResponseWriter, req *http.Request) {
	session := data.GetOrCreateSession(req)
	log.Printf("%s [%s]", req.URL, session.Id)
	SendPage(w, req, session, "about", "Acerca de")
}


var port = flag.Int("port", 8080, "Network port")


func Listen() {
	p := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(p, nil)
	if err != nil {
		log.Fatalf("Cannot Listen: %s", err)
	}
}

func main() {
	flag.Parse()

	mime.AddExtensionType(".ttf", "font/ttf")

	content.WatchForChanges(func(id string) {
		if id == "" {
			cache.Touch("/courses")
		} else {
			cache.Touch("/content/" + id)
		}
	})

	// handlers
	http.HandleFunc("/css/", StaticFiles)
	http.HandleFunc("/js/", StaticFiles)
	http.HandleFunc("/img/", StaticFiles)
	http.HandleFunc("/fonts/", StaticFiles)

	http.HandleFunc("/favicon.ico", hFavicon)
	http.HandleFunc("/png/", hPhotos)

	/*

		// Users
		http.HandleFunc("/login", hLogin)
		http.HandleFunc("/logout", hLogout)
		http.HandleFunc("/register", hRegister)

		// About
		http.HandleFunc("/acerca", hAbout)

	*/

	http.HandleFunc("/", fragmentPage)

	http.HandleFunc("/_frag/", hFragList)

	Listen()
}
