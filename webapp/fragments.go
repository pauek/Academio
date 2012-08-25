package main

import (
	"Academio/content"
	"Academio/webapp/data"
	"bytes"
	"encoding/json"
	"html/template"
	"fmt"
	F "fragments"
	"log"
	"net/http"
	"strings"
	"time"
)

func fragmentPage(w http.ResponseWriter, req *http.Request) {
	// Determine fragment + title
	var title, fid string
	title, fid, notfound := pathToFragmentID(req.URL.Path[1:])
	if notfound {
		log.Printf("%s (NOT FOUND)", req.URL)
		http.NotFound(w, req)
		return
	}

	// Get session
	session := data.GetOrCreateSession(req)
	id := session.Id
	if session.User != nil {
		id = session.User.Login
	}
	log.Printf("%s [%s]", req.URL, id)
	session.PutCookie(w)

	ModeDispatch(w, req, session, fid, title)
}

func ModeDispatch(w http.ResponseWriter, req *http.Request, session *data.Session, fid, title string) {
	// Send HTML or JSON
	switch req.Header.Get("Fragments") {
	case "":
		sendHTML(w, session, fid, title)
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

type layoutInfo struct {
	Title string
	Message string
	Navbar template.HTML
	Body template.HTML
}

func sendHTML(w http.ResponseWriter, session *data.Session, fid, title string) {
	w.Header().Set("Content-Type", "text/html")
	navbarfid := "navbar"
	if session.User != nil {
		navbarfid += " " + session.Id
	}
	if layout := tmpl.Lookup("layout"); layout != nil {
		layout.Execute(w, layoutInfo{
			Title: title,
			Message: session.Message,
			Navbar: 	template.HTML(cache.RenderToString(navbarfid)),
			Body: template.HTML(cache.RenderToString(fid)),
		})
	} else {
		code := http.StatusInternalServerError
		http.Error(w, "Template 'layout' not found", code)
	}
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
	cache.Register(fNavbar, "navbar")
	cache.Register(fItem, "item")
	cache.Register(fItemFragment,
		"item-nav",
		"item-link",
		"topic-small",
		"concept-small",
	)
	cache.Register(fStatic,
		"home",
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

type navbarInfo struct {
	User    *data.User
	Message *string
}
func fNavbar(c *F.Cache, args []string) F.Fragment {
	fmt.Printf("fNavbar: args = %v\n", args)
	fid := strings.Join(args, " ")
	c.Depends(fid, "/templates")
	var info navbarInfo
	if len(args) > 1 {
		session := data.FindSession(args[1])
		info.User = session.User
		c.Depends(fid, "/session/" + session.Id)
		if len(session.Message) > 0 {
			info.Message = &session.Message
		}
	}
	return F.MustParse(exec(args[0], info))
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
		err := t.Execute(&b, data)
		if err != nil {
			return err.Error()
		} else {
			return b.String()
		}
	}
	panic("missing template")
}
