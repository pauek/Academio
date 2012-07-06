package main

import (
	ct "Academio/content"
	"bytes"
	"fmt"
	frag "fragments"
	T "html/template"
	_ "log"
	"net/http"
	_ "strings"
)

var (
	tmpl   = T.Must(T.ParseGlob("templates/" + "*.html"))
	layout = frag.MustParseFile("templates/layout")
)

func exec(tname string, data interface{}) frag.Fragment {
	var b bytes.Buffer
	if t := tmpl.Lookup(tname); t != nil {
		t.Execute(&b, data)
		return frag.Text(b.String())
	}
	panic("missing template")
}

func fItem(C *frag.Cache, args []string) frag.Fragment {
	item := ct.Get(args[1])
	return exec(item.Type(), item)
}

func fStatic(C *frag.Cache, args []string) frag.Fragment {
	return exec(args[0], nil)
}

func fCourses(C *frag.Cache, args []string) frag.Fragment {
	return exec("courses", ct.Courses())
}

func hRoot(w http.ResponseWriter, req *http.Request) {
	var title, fid string // fid = fragment id
	path := req.URL.Path[1:]
	switch path {
	case "":
		title, fid = "Inicio", "home"
	case "cursos":
		title, fid = "Cursos", "courses"
	default:
		item := ct.Get(path)
		if item == nil {
			http.NotFound(w, req)
			return
		}
		title = item.Data().Title
		fid = fmt.Sprintf("item %s", path)
	}
	layout.Exec(w, func(action string) {
		switch action {
		case "body":
			frag.Render(w, fid)
		case "title":
			fmt.Fprintf(w, title)
		default:
			frag.Render(w, action)
		}
	})
}

func main() {
	ct.WatchForChanges(func(id string) {
		if id == "" {
			frag.Invalidate("courses")
		} else {
			frag.Invalidate("item " + id)
		}
	})

	// fragments
	frag.Register("item", fItem)
	frag.Register("home", fStatic)
	frag.Register("navbar", fStatic)
	frag.Register("footer", fStatic)
	frag.Register("courses", fCourses)

	// handlers
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", hRoot)

	http.ListenAndServe(":8080", nil)
}
