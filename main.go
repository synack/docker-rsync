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
		fmt.Fprintf(os.Stderr, "   or: %v [options] rsync://IP:PORT/MODULE\n", os.Args[0])
		fmt.Print("\nOptions:\n")
		flag.PrintDefaults()
	}

	var version = flag.Bool("version", false, "Print version")
	var watch = flag.Bool("watch", true, "Watch source directory for changes")
	var verbose = flag.Bool("verbose", false, "Verbose output")
	var srcpath = flag.String("src", pwd, "Source directory")
	var dstpath = flag.String("dst", pathpkg.Join("/rsync", pwd), "Destination directory")

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

	if *srcpath == "" {
		fmt.Println("error: please specify -src argument")
		os.Exit(1)
	}

	if *dstpath == "" {
		fmt.Println("error: please specify -dst argument")
		os.Exit(1)
	}

	rpath := *srcpath
	rpathDir := pathpkg.Dir(*dstpath)

	// TODO: refactor the following part...

	if strings.HasPrefix(via, "rsync://") {
		// use rsync protocol directly
		rsyncEndpoint := via

		fmt.Printf("Syncing %s (local) to %s (%s)\n", *srcpath, *dstpath, rsyncEndpoint)
		Sync(rsyncEndpoint, 0, rpath, rpathDir, *verbose) // initial sync

		if *watch {
			fmt.Println("Watching for file changes ...")
			Watch(rpath, func(id uint64, path string, flags []string) {
				Sync(rsyncEndpoint, 0, rpath, rpathDir, true)
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

		Provision(machineName, *verbose)
		RunSSHCommand(machineName, "sudo mkdir -p "+rpathDir, *verbose)
		fmt.Printf("Syncing %s (local) to %s (docker-machine %s)\n", *srcpath, *dstpath, machineName)
		Sync(machineName, port, rpath, rpathDir, *verbose) // initial sync

		if *watch {
			fmt.Println("Watching for file changes ...")
			Watch(rpath, func(id uint64, path string, flags []string) {
				Sync(machineName, port, rpath, rpathDir, true)
			})
		}
	}

}
