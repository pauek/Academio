package main

import (
	"Academio/content"
	"Academio/webapp/data"
	"bytes"
	"encoding/json"
	"fmt"
	F "fragments"
	"log"
	"net/http"
	"time"
)

func fragmentPage(w http.ResponseWriter, req *http.Request) {
	session := data.GetSession(w, req)
	id := session.Id
	if session.User != nil {
		id = session.User.Login
	}
	log.Printf("%s [%s]", req.URL, id)

	// Determine fragment + title
	var title, fid string
	title, fid, notfound := pathToFragmentID(req.URL.Path[1:])
	if notfound {
		http.NotFound(w, req)
		return
	}

	// Send HTML or JSON
	switch req.Header.Get("Fragments") {
	case "":
		sendHTML(w, title, fid)
	case "all":
		sendJSON(w, cache.List(fid))
	case "since":
		sendJSON(w, cache.ListDiff(fid, getFragmentsStamp(req)))
	}
}

func pathToFragmentID(path string) (title, fid string, notfound bool) {
	switch path {
	case "":
		title, fid = "Inicio", "home"
	case "cursos":
		title, fid = "Cursos", "courses"
	default:
		item := content.Get(path)
		if item == nil {
			return "", "", true
		}
		title = item.Data().Title
		fid = fmt.Sprintf("item %s", path)
	}
	return
}

func getFragmentsStamp(req *http.Request) (stamp time.Time) {
	since := req.Header.Get("FragmentsStamp")
	err := json.Unmarshal([]byte(since), &stamp)
	if err != nil {
		log.Printf("ERROR: Cannot unmarshal timestamp '%s'", since)
	}
	return
}

func sendHTML(w http.ResponseWriter, title, fid string) {
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
}

func sendJSON(w http.ResponseWriter, list []F.ListItem) {
	if data, err := json.Marshal(list); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		log.Printf("ERROR: Cannot marshal fragment: %s", err)
	}
}

// Fragments

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
		"login",
	)
	cache.Register(fCourses, "courses")
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

func exec(tname string, data interface{}) string {
	var b bytes.Buffer
	if t := tmpl.Lookup(tname); t != nil {
		t.Execute(&b, data)
		return b.String()
	}
	panic("missing template")
}
