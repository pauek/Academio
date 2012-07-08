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

var (
	tmpl   = T.Must(T.New("").Funcs(tFuncs).ParseGlob("templates/" + "[a-zA-Z0-9]*.html"))
	layout = F.MustParse(exec("layout", nil))
	cache  = F.NewCache()
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
	C.Depends("item "+args[1],
		"/item/"+args[1],
		"/templates/"+item.Type()+".html",
	)
	return F.MustParse(exec(item.Type(), item))
}

func fItemFragment(C *F.Cache, args []string) F.Fragment {
	C.Depends(args[0]+" "+args[1], 
		"/item/"+args[1],
		"/templates/"+args[0]+".html",
	)
	return F.MustParse(exec(args[0], content.Get(args[1])))
}

func fStatic(C *F.Cache, args []string) F.Fragment {
	C.Depends(args[0], "/templates/static.html")
	return F.Text(exec(args[0], nil))
}

func fCourses(C *F.Cache, args []string) F.Fragment {
	C.Depends("courses", "/courses")
	return F.Text(exec("courses", content.Courses()))
}

func hFragList(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/_frag/"):]
	list := cache.List("item " + id)
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
				cache.Render(w, fid)
			case "title":
				fmt.Fprintf(w, title)
			default:
				cache.Render(w, action)
			}
		})
	case "all":
		list := cache.List(fid)
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
		list := cache.ListDiff(fid, stamp)
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
			cache.Touch("/courses")
		} else {
			cache.Touch("/item/" + id)
		}
	})
	watchTemplates()

	// fragments
	cache.Register(fItem, "item")
	cache.Register(fItemFragment,
		"item-nav",
		"item-link",
		"topic-small",
		"concept-small",
	)
	cache.Register(fStatic,
		"home",
		"navbar",
		"footer",
	)
	cache.Register(fCourses, "courses")

	// handlers
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/_frag/", hFragList)
	http.HandleFunc("/", Page)

	http.ListenAndServe(":8080", nil)
}
