package main

import (
	"text/template"
	"io"
	"log"
	"net"
	"net/http"
)

const webapp = "localhost:8080"

type Server string

func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	var inConn, outConn net.Conn

	log.Printf("%s %v", req.Method, req.URL)

	if inConn, err = net.Dial("tcp", string(s)); err != nil {
		SiteIsDown(w)
		return
	}
	if outConn, _, err = w.(http.Hijacker).Hijack(); err != nil {
		// log.Printf("Cannot hijack connection: %s", err)
		inConn.Close()
		return
	}

	go func() {
		io.Copy(outConn, inConn)
		outConn.Close()
	}()
	go func() {
		req.Header["X-Forwarded-For"] = []string{req.RemoteAddr}
		req.Write(inConn)
		io.Copy(inConn, outConn)
		inConn.Close()
	}()
	return
}

const down = `<!doctype html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>Academio duerme</title>
    <style>
      html, body { 
        margin: 0; padding: 0;
        width: 100%; 
        background: #29303d;
      }
      h1 {
        color: white;
        font-weight: normal;
        text-align: center;
        position: absolute;
        top: 50%;
        width: 100%;
        margin-top: -1em;
        font-family: "Open Sans", sans-serif;
        font-size: 2em;
        line-height: 1em;
      }
      .small {
        color: #abb;
        font-weight: normal;
        font-size: .6em;
      }
      .orange {
        color: orange;
      }
      .bold { font-weight: bold; }
    </style>
  </head>
  <body>{{.}}</body>
</html>
`

const downBody = `
<h1>
  <span class="bold"><span class="orange">academ.</span>io</span> duerme <br />
  <span class="small">Por favor, vuelve más tarde</span>
</h1>
`

var tDown *template.Template

func init() {
	var err error
	if tDown, err = template.New("").Parse(down); err != nil {
		log.Fatalf("Template 'down' error: %s", err)
	}
}

func SiteIsDown(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	tDown.Execute(w, downBody)
}

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("Cannot listen on ':80': %s", err)
	}
	http.Serve(listener, Server(webapp))
}
