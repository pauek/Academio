package content

import (
	"inotify"
	"log"
	"path/filepath"
)

func isDir(ev *inotify.Event) bool {
	return (ev.Mask & inotify.IN_ISDIR) != 0
}

func isChange(ev *inotify.Event) bool {
	ch := inotify.IN_CREATE
	ch |= inotify.IN_DELETE
	ch |= inotify.IN_MODIFY
	ch |= inotify.IN_MOVE
	return (ev.Mask & ch) != 0
}

func WatchForChanges(onChange func(id string)) {
	// Create watcher
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Printf("warning: Cannot create inotify.Watcher")
		return
	}

	// watch dirs in levels 1-3 from the roots
	watchList := []string{} // prevent walkDirs from generating events
	for _, root := range roots {
		walkDirs(root, func(reldir string, level int) {
			watchList = append(watchList, filepath.Join(root, reldir))
		})
	}
	for _, dir := range watchList {
		watcher.Watch(dir)
	}

	// go wait for events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				rel := removeRoot(ev.Name)
				if onChange != nil && isChange(ev) {
					if !isDir(ev) {
						rel = filepath.Dir(rel)
					}
					id := toID(rel)
					onChange(id)
					if isDir(ev) {
						onChange(parentID(id))
					}
				}
				// watch new dirs
				if isDir(ev) {
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
