package main

import (
	"fmt"
	F "fragments"
	T "html/template"
	"inotify"
	"log"
)

var tFuncs = map[string]interface{}{
	"plus1": plus1,
}

func plus1(i int) int {
	return i + 1
}

var (
	tmpl   *T.Template
	layout F.Template
)

func init() {
	readTemplates();
	watchTemplates();
}

func readTemplates() {
	tmpldir := srvdir + "/templates/[a-zA-Z0-9]*.html"
	tmpl = T.Must(T.New("").Funcs(tFuncs).ParseGlob(tmpldir))
	layout = F.MustParse(exec("layout", nil))
}

func isChange(ev *inotify.Event) bool {
	ch := inotify.IN_CREATE
	ch |= inotify.IN_DELETE
	ch |= inotify.IN_MODIFY
	ch |= inotify.IN_MOVE
	return (ev.Mask & ch) != 0
}

func watchTemplates() {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Printf("Warning: Cannot watch templates.")
	}

	watcher.Watch("templates")

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isChange(ev) {
					fmt.Printf("Changed: %s\n", ev.Name)
					readTemplates()
					cache.Touch("/templates")
				}
			}
		}
	}()
}

