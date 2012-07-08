package main

import (
	"fmt"
	"inotify"
	"log"
)

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