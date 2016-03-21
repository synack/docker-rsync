package main

import (
	"github.com/fsnotify/fsevents"
	"sort"
	"time"
)

var noteDescription = map[fsevents.EventFlags]string{
	fsevents.MustScanSubDirs: "MustScanSubdirs",
	fsevents.UserDropped:     "UserDropped",
	fsevents.KernelDropped:   "KernelDropped",
	fsevents.EventIDsWrapped: "EventIDsWrapped",
	fsevents.HistoryDone:     "HistoryDone",
	fsevents.RootChanged:     "RootChanged",
	fsevents.Mount:           "Mount",
	fsevents.Unmount:         "Unmount",

	fsevents.ItemCreated:       "Created",
	fsevents.ItemRemoved:       "Removed",
	fsevents.ItemInodeMetaMod:  "InodeMetaMod",
	fsevents.ItemRenamed:       "Renamed",
	fsevents.ItemModified:      "Modified",
	fsevents.ItemFinderInfoMod: "FinderInfoMod",
	fsevents.ItemChangeOwner:   "ChangeOwner",
	fsevents.ItemXattrMod:      "XAttrMod",
	fsevents.ItemIsFile:        "IsFile",
	fsevents.ItemIsDir:         "IsDir",
	fsevents.ItemIsSymlink:     "IsSymLink",
}

func Watch(path string, eventHandler func(id uint64, path string, flags []string)) {
	dev, _ := fsevents.DeviceForPath(path)
	fsevents.EventIDForDeviceBeforeTime(dev, time.Now())

	es := &fsevents.EventStream{
		Paths:   []string{path},
		Latency: 50 * time.Millisecond,
		Device:  dev,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot}
	es.Start()
	ec := es.Events

	for {
		select {
		case event := <-ec:
			flags := make([]string, 0)
			for bit, description := range noteDescription {
				if event.Flags&bit == bit {
					flags = append(flags, description)
				}
			}
			sort.Sort(sort.StringSlice(flags))
			go eventHandler(event.ID, event.Path, flags)

			es.Flush(false)
		}
	}
}
