package main

import (
	"Academio/content"
	"bytes"
	"encoding/json"
	"fmt"
	F "fragments"
	"log"
	"net/http"
	"time"
)

func init() {
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
}

func fragmentPage(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s", req.URL)
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
		w.Header().Set("Content-Type", "text/html")
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
		"/content/"+args[1],
		"/templates",
	)
	if topic, ok := item.(*content.Topic); ok {
		for _, subitems := range topic.Children() {
			C.Depends("item "+args[1], "/content/"+subitems.Id)
		}
	}
	return F.MustParse(exec(item.Type(), item))
}

func fItemFragment(C *F.Cache, args []string) F.Fragment {
	C.Depends(args[0]+" "+args[1],
		"/content/"+args[1],
		"/templates",
	)
	return F.MustParse(exec(args[0], content.Get(args[1])))
}

func fStatic(C *F.Cache, args []string) F.Fragment {
	C.Depends(args[0], "/templates")
	return F.Text(exec(args[0], nil))
}

func fCourses(C *F.Cache, args []string) F.Fragment {
	C.Depends("courses",
		"/courses",
		"/templates",
	)
	return F.MustParse(exec("courses", content.Courses()))
}
