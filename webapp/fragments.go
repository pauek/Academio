package main

import (
	"Academio/content"
	"Academio/webapp/data"
	"bytes"
	"encoding/json"
	"fmt"
	F "fragments"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func fragmentPage(w http.ResponseWriter, req *http.Request) {

	// Do not show main page for now
	if req.URL.Path == "/" {
		http.Redirect(w, req, "/cursos", http.StatusSeeOther)
		return
	}

	// Extract Item ID from URL
	var itemId string
	path := req.URL.Path[1:]
	k := strings.Index(path, "/")
	if k != -1 {
		// serve file from item directory
		itemId = path[:k]
		dir := content.ToDir(itemId)
		if dir.Root == "" {
			NotFound(w, req)
			return
		} 
		filename := dir.File(path[k+1:])
		log.Printf("%s", req.URL.Path)
		http.ServeFile(w, req, filename)
		return
	} else {
		itemId = path
	}

	// Determine fragment ID & Title
	var fid string
	fid, notfound := getFragmentID(itemId)
	if notfound {
		log.Printf("%s (NOT FOUND)", req.URL)
		NotFound(w, req)
		return
	}
	title := getTitle(itemId)

	// Get session
	session := data.GetOrCreateSession(req)
	id := session.Id
	if session.User != nil {
		id = session.User.Login
	}
	log.Printf("%s [%s]", req.URL, id)
	session.PutCookie(w)

	SendPage(w, req, session, fid, title)
	session.Message = "" // message delivered
}

func SendPage(w http.ResponseWriter, req *http.Request, session *data.Session, fid, title string) {
	_stamp := req.Header.Get("FragmentsSince")
	if _stamp == "" {
		sendHTML(w, session, fid, title)
		return
	}
	stamp := parseStamp(_stamp)
	sendJSON(w, &JSONPage{
		Stamp:   time.Now(),
		Title:   title,
		Message: session.Message,
		Navbar:  cache.Diff(navbarFID(session), stamp),
		Body:    cache.Diff(fid, stamp),
	})
}

func getFragmentID(id string) (fid string, notfound bool) {
	switch id {
	case "":
		fid = "home"
	case "cursos":
		fid = "courses"
	default:
		item := content.Get(id)
		if item == nil {
			return "", true
		}
		fid = fmt.Sprintf("item %s", id)
	}
	return
}

func getTitle(id string) (title string) {
	switch id {
	case "":
		title = "Inicio"
	case "cursos":
		title = "Cursos"
	default:
		title = content.Get(id).Data().Title
	}
	return
}

func navbarFID(session *data.Session) string {
	navbarfid := "navbar"
	if session != nil && session.User != nil {
		navbarfid += " " + session.Id
	}
	return navbarfid
}

func parseStamp(s string) (stamp time.Time) {
	if s == "null" {
		return
	}
	err := json.Unmarshal([]byte(s), &stamp)
	if err != nil {
		log.Printf("ERROR: Cannot unmarshal timestamp '%s'", s)
	}
	return
}

type layoutInfo struct {
	Title   string
	Message string
	Navbar  template.HTML
	Body    template.HTML
}

func sendHTML(w http.ResponseWriter, session *data.Session, fid, title string) {
	w.Header().Set("Content-Type", "text/html")
	if layout := tmpl.Lookup("layout"); layout != nil {
		var msg string
		if session != nil {
			msg = session.Message
		}
		layout.Execute(w, layoutInfo{
			Title:   title,
			Message: msg,
			Navbar:  template.HTML(cache.RenderToString(navbarFID(session))),
			Body:    template.HTML(cache.RenderToString(fid)),
		})
	} else {
		code := http.StatusInternalServerError
		http.Error(w, "Template 'layout' not found", code)
	}
}

type JSONPage struct {
	Stamp   time.Time
	Title   string
	Message string
	Navbar  []F.ListItem
	Body    []F.ListItem
}

func sendJSON(w http.ResponseWriter, page *JSONPage) {
	if data, err := json.Marshal(page); err == nil {
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
		"about",
		"notfound",
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
		c.Depends(fid, "/session/"+session.Id)
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
