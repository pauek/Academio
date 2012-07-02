package main

import (
	ct "Academio/content"
	"bytes"
	"fmt"
	frag "fragments"
	"io"
	"net/http"
	T "text/template"
)

const tLayout = `<!doctype html>
<html>
  <head>
  </head>
  <body>{{}}</body>
</html>`

const tCourse = `<h1>{{.Title}}</h1>
<p>{{.Html}}</p>
<ul>
{{range .Topics}}
<li><a href="/{{.}}">{{.}}</a></li>
{{end}}
</ul>`

const tTopic = `<h1>{{.Title}}</h1>
<p>{{.Html}}</p>
<ul>
{{range .Concepts}}
<li><a href="/{{.}}">{{.}}</a></li>
{{end}}
</ul>`

const tConcept = `<h1>{{.Title}}</h1>
<p><b>Video:</b> {{.VideoID}}</p>
<p>{{.Html}}</p>`

var (
	layout  frag.Template
	course  *T.Template
	topic   *T.Template
	concept *T.Template
)

func init() {
	layout, _ = frag.Parse(tLayout)
	course, _ = T.New("").Parse(tCourse)
	topic, _ = T.New("").Parse(tTopic)
	concept, _ = T.New("").Parse(tConcept)
}

func fCourse(C *frag.Cache, args []string) frag.Fragment {
	var b bytes.Buffer
	course.Execute(&b, ct.GetCourse(ct.Id(args[1])))
	return frag.Text(b.String())
}

func fTopic(C *frag.Cache, args []string) frag.Fragment {
	var b bytes.Buffer
	topic.Execute(&b, ct.GetTopic(ct.Id(args[1])))
	return frag.Text(b.String())
}

func fConcept(C *frag.Cache, args []string) frag.Fragment {
	var b bytes.Buffer
	c := ct.GetConcept(ct.Id(args[1]))
	fmt.Sprintf("%s\n", c.Doc.Html)
	concept.Execute(&b, c)
	return frag.Text(b.String())
}

func Page(ref string) frag.Fragment {
	return layout.RenderFn(func(w io.Writer, id string) {
		frag.Render(w, ref)
	})
}

func fPage(C *frag.Cache, args []string) frag.Fragment {
	switch ct.Type(ct.Id(args[1])) {
	case ct.CourseType:
		return Page("course " + args[1])
	case ct.TopicType:
		return Page("topic " + args[1])
	case ct.ConceptType:
		return Page("concept " + args[1])
	}
	return frag.Text("Error page")
}

func hRoot(w http.ResponseWriter, req *http.Request) {
	frag.Render(w, "page "+req.URL.Path[1:])
}

func main() {
	frag.Register("course", fCourse)
	frag.Register("topic", fTopic)
	frag.Register("concept", fConcept)
	frag.Register("page", fPage)
	http.HandleFunc("/", hRoot)
	http.ListenAndServe(":8080", nil)
}
