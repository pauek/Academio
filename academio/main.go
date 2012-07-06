package main

import (
	ct "Academio/content"
	"bytes"
	"fmt"
	frag "fragments"
	T "html/template"
	"net/http"
)

var tFuncs = map[string]interface{}{
	"plus1": func(i int) int { return i + 1 },
}

var (
	tmpl   = T.Must(T.New("").Funcs(tFuncs).ParseGlob("templates/" + "*.html"))
	layout = frag.MustParseFile("templates/layout")
)

func exec(tname string, data interface{}) string {
	var b bytes.Buffer
	if t := tmpl.Lookup(tname); t != nil {
		t.Execute(&b, data)
		return b.String()
	}
	panic("missing template")
}

func fItem(C *frag.Cache, args []string) frag.Fragment {
	item := ct.Get(args[1])
	return frag.MustParse(exec(item.Type(), item))
}

func fTopicSmall(C *frag.Cache, args []string) frag.Fragment {
	topic := ct.Get(args[1])
	return frag.Text(exec("topic-small", topic))
}

func fStatic(C *frag.Cache, args []string) frag.Fragment {
	return frag.Text(exec(args[0], nil))
}

func fCourses(C *frag.Cache, args []string) frag.Fragment {
	return frag.Text(exec("courses", ct.Courses()))
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
	frag.Register("topic-small", fTopicSmall)

	// handlers
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", hRoot)

	http.ListenAndServe(":8080", nil)
}
