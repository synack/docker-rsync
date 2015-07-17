package main

import (
	"golang.org/x/exp/winfsnotify" "wfsevent"
)

func Watch(path string, eventHandler func(id uint64, path string, flags []string)) {

	
	es := &wfsevent.NewWatcher.Watch{[]string{path}}
	ec := es.Event.String

	for {
		select {
		case event := <-ec:
			go eventHandler(1 uint64, "dummy", [0])
		}
	}
}
