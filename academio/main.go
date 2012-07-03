package main

import (
	ct "Academio/content"
	"bytes"
	"fmt"
	frag "fragments"
	T "html/template"
	"net/http"
)

const tLayout = `<!doctype html>
<html>
  <head></head>
  <body>{{}}</body>
</html>`

var (
	tmpl   = T.Must(T.ParseFiles("template.html"))
	layout = frag.MustParse(tLayout)
)

func fItem(C *frag.Cache, args []string) frag.Fragment {
	var b bytes.Buffer
	item := ct.Get(args[1])
	tmpl.Lookup(ct.Type(item)).Execute(&b, item)
	return frag.Text(b.String())
}

func fCourseList(C *frag.Cache, args []string) frag.Fragment {
	var b bytes.Buffer
	tmpl.Lookup("courselist").Execute(&b, ct.Courses)
	return frag.Text(b.String())
}

func hRoot(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[1:]
	if len(id) > 0 {
		item := ct.Get(id)
		if item == nil {
			http.NotFound(w, req)
			return
		}
		id = fmt.Sprintf("item %s", id)
	} else {
		id = "courselist"
	}
	layout.Exec(w, func(action string) {
		frag.Render(w, id)
	})
}

func main() {
	ct.Read()
	frag.Register("item", fItem)
	frag.Register("courselist", fCourseList)
	http.HandleFunc("/", hRoot)
	http.ListenAndServe(":8080", nil)
}
