package main

import (
	ct "Academio/content"
	"bytes"
	"fmt"
	frag "fragments"
	T "html/template"
	"net/http"
	"strings"
)

const tLayout = `<!doctype html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
  <title>Academio - {{title}}</title>
  <link rel="stylesheet/less" type="text/css" href="/css/academio.less">
  <script src="/js/less-1.3.0.min.js" type="text/javascript"></script>
</head>
<body>
  <header>{{navbar}}</header>
  <div id="content">{{body}}</div>
  <footer>{{footer}}</footer>
</body>
</html>`

var (
	tmpl   = T.Must(T.ParseFiles("template.html"))
	layout = frag.MustParse(tLayout)
)

func exec(tname string, data interface{}) frag.Fragment {
	var b bytes.Buffer
	tmpl.Lookup(tname).Execute(&b, data)
	return frag.Text(b.String())
}

func fItem(C *frag.Cache, args []string) frag.Fragment {
	var b bytes.Buffer
	item := ct.Get(args[1])
	tmpl.Lookup(item.Type()).Execute(&b, item)
	C.Invalidate(strings.Join(args, " ")) // always invalid
	return frag.Text(b.String())
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
