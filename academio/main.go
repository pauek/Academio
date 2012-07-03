package main

import (
	ct "Academio/content"
	"bytes"
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
	tmpl = T.Must(T.ParseFiles("template.html"))
	layout = frag.MustParse(tLayout)
)

func fItem(C *frag.Cache, args []string) frag.Fragment {
	var b bytes.Buffer
	item := ct.Get(args[1])
	tmpl.Lookup(ct.Type(item)).Execute(&b, item)
	return frag.Text(b.String())
}

func hRoot(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[1:]
	item := ct.Get(id)
	if item == nil {
		http.NotFound(w, req)
		return
	}
	layout.Exec(w, func(action string) {
		frag.Render(w, "item "+id)
	})
}

func main() {
	ct.Read()
	ct.Show()
	frag.Register("item", fItem)
	http.HandleFunc("/", hRoot)
	http.ListenAndServe(":8080", nil)
}
