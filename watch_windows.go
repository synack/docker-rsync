// +build windows

package main

import (
	"fmt"
	"golang.org/x/exp/winfsnotify"
)

var noteDescription = map[uint32]string{
	winfsnotify.FS_CREATE: "Created",
	winfsnotify.FS_DELETE: "Removed",
	winfsnotify.FS_MODIFY: "Modified",
	winfsnotify.FS_MOVE:   "Renamed",
}

func Watch(path string, eventHandler func(id uint64, path string, flags []string)) {
	w, err := winfsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("Cannot start watcher: %s", err))
	}

	err = w.Watch(path)
	if err != nil {
		panic(fmt.Sprintf("Cannot watch path: %s", err))
	}

	for {
		select {
		case event := <-w.Event:
			fmt.Printf("Got event %#v\n\n", event)

			if mask, ok := noteDescription[event.Mask]; ok {
				go eventHandler(uint64(event.Mask), event.Name, []string{mask})
			}
		}
	}
}
