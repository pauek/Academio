package content

import (
	"inotify"
	"log"
	"path/filepath"
)

var OnChange func(id string)

func WatchForChanges() {
	// Create watcher
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Printf("warning: Cannot create inofity.Watcher")
		return
	}

	// watch dirs in levels 1-3 from the roots
	var path string
	var walk func(dir string, level int) // recursive

	walk = func(dir string, level int) {
		absdir := filepath.Join(path, dir)
		if level > 0 {
			// fmt.Printf("%s\n", absdir)
			watcher.Watch(absdir)
		}
		eachSubDir(absdir, func(subdir string) {
			reldir := filepath.Join(dir, subdir)
			if level < 3 {
				walk(reldir, level+1)
			}
		})
	}
	for _, path = range roots {
		walk("", 0)
	}

	// go wait for events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				rel := removeRoot(ev.Name)
				if OnChange != nil {
					OnChange(toID(rel))
				}
				// watch new dirs
				if ev.Mask & inotify.IN_ISDIR != 0 {
					lv := numLevels(rel)
					if lv > 0 && lv < 4 {
						watcher.Watch(ev.Name)
					}
				}
			case err := <-watcher.Error:
				log.Printf("watcher error: %s", err)
			}
		}
	}()
}
