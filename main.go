package main

import (
	"flag"
	"fmt"
	"os"
	pathpkg "path"
	"strings"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error: unable to get current directory:", err)
		os.Exit(1)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %v [options] DOCKER-MACHINE-NAME\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "or   : %v [options] rsync://IP:PORT/MODULE\n", os.Args[0])
		fmt.Print("\nOptions:\n")
		flag.PrintDefaults()
	}

	var version = flag.Bool("version", false, "Print version")
	var onetime = flag.Bool("1", false, "Sync only once")
	var path = flag.String("path", pwd, "Sync this directory")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	via := flag.Args()[0]

	rpath := *path
	rpathDir := pathpkg.Dir(*path)

	// TODO: refactor the following part...

	if strings.HasPrefix(via, "rsync://") {
		// use rsync protocol directly
		rsyncEndpoint := via

		Sync(rsyncEndpoint, 0, rpath, rpathDir) // initial sync

		if !*onetime {
			Watch(rpath, func(id uint64, path string, flags []string) {
				Sync(rsyncEndpoint, 0, rpath, rpathDir)
			})
		}

	} else {
		// use rsync via ssh
		machineName := via

		port, err := GetSSHPort(machineName)
		if err != nil {
			fmt.Printf("error: unable to get port for machine '%v': %v\n", machineName, err)
			os.Exit(1)
		}

		Provision(machineName)
		RunSSHCommand(machineName, "sudo mkdir -p "+rpathDir)
		Sync(machineName, port, rpath, rpathDir) // initial sync

		if !*onetime {
			Watch(rpath, func(id uint64, path string, flags []string) {
				Sync(machineName, port, rpath, rpathDir)
			})
		}
	}

}
