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
		log.Printf("warning: Cannot create inotify.Watcher")
		return
	}

	// watch dirs in levels 1-3 from the roots
	for _, root := range roots {
		walkDirs(root, func(reldir string, level int) {
			watcher.Watch(filepath.Join(root, reldir))
		})
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
