package main

import (
	"Academio/content"
	"bytes"
	"encoding/json"
	"fmt"
	F "fragments"
	T "html/template"
	"log"
	"net/http"
	"time"
)

var tFuncs = map[string]interface{}{
	"plus1": func(i int) int { return i + 1 },
}

var (
	tmpl   = T.Must(T.New("").Funcs(tFuncs).ParseGlob("templates/" + "*.html"))
	layout = F.MustParseFile("templates/layout")
)

func exec(tname string, data interface{}) string {
	var b bytes.Buffer
	if t := tmpl.Lookup(tname); t != nil {
		t.Execute(&b, data)
		return b.String()
	}
	panic("missing template")
}

func fItem(C *F.Cache, args []string) F.Fragment {
	item := content.Get(args[1])
	return F.MustParse(exec(item.Type(), item))
}

func fItemFragment(C *F.Cache, args []string) F.Fragment {
	return F.MustParse(exec(args[0], content.Get(args[1])))
}

func fStatic(C *F.Cache, args []string) F.Fragment {
	return F.Text(exec(args[0], nil))
}

func fCourses(C *F.Cache, args []string) F.Fragment {
	return F.Text(exec("courses", content.Courses()))
}

func hFragList(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/_frag/"):]
	list := F.List("item " + id)
	tmpl.Lookup("fraglist").Execute(w, list)
}

func Page(w http.ResponseWriter, req *http.Request) {
	var title, fid string // fid = fragment id
	path := req.URL.Path[1:]
	switch path {
	case "":
		title, fid = "Inicio", "home"
	case "cursos":
		title, fid = "Cursos", "courses"
	default:
		item := content.Get(path)
		if item == nil {
			http.NotFound(w, req)
			return
		}
		title = item.Data().Title
		fid = fmt.Sprintf("item %s", path)
	}
	switch req.Header.Get("Fragments") {
	case "":
		layout.Exec(w, func(action string) {
			switch action {
			case "body":
				F.Render(w, fid)
			case "title":
				fmt.Fprintf(w, title)
			default:
				F.Render(w, action)
			}
		})
	case "all":
		list := F.List(fid)
		if data, err := json.Marshal(list); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else {
			log.Printf("ERROR: Cannot marshal fragment: %s", err)
		}
	case "since":
		since := req.Header.Get("FragmentsStamp")
		var stamp time.Time
		err := json.Unmarshal([]byte(since), &stamp)
		if err != nil {
			log.Printf("ERROR: Cannot unmarshal timestamp '%s'", since)
		}
		list := F.ListDiff(fid, stamp)
		if data, err := json.Marshal(list); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else {
			log.Printf("ERROR: Cannot marshal fragment: %s", err)
		}
	}
}

func main() {
	content.WatchForChanges(func(id string) {
		if id == "" {
			F.Invalidate("courses")
		} else {
			F.Invalidate("item " + id)
		}
	})

	// fragments
	F.Register("item", fItem)
	F.Register("item-nav", fItemFragment)
	F.Register("item-link", fItemFragment)
	F.Register("topic-small", fItemFragment)
	F.Register("concept-small", fItemFragment)
	F.Register("home", fStatic)
	F.Register("navbar", fStatic)
	F.Register("footer", fStatic)
	F.Register("courses", fCourses)

	// handlers
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/_frag/", hFragList)
	http.HandleFunc("/", Page)

	http.ListenAndServe(":8080", nil)
}
