package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	pathpkg "path"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %v [-version|-initial] <machine-name>\n", os.Args[0])
		flag.PrintDefaults()
	}

	var version = flag.Bool("version", false, "Print version")
	var onetime = flag.Bool("1", false, "Sync only once")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("error: unable to get current directory:", err)
		os.Exit(1)
	}

	if len(flag.Args()) != 1 {
		fmt.Printf("usage: %v [-version|-initial] <machine-name>\n", os.Args[0])
		os.Exit(1)
	}

	machineName := flag.Args()[0]

	port, err := GetSSHPort(machineName)
	if err != nil {
		fmt.Printf("error: unable to get port for machine '%v': %v\n", machineName, err)
		os.Exit(1)
	}

	Provision(machineName)

	rpath := path
	rpathDir := pathpkg.Dir(path)

	PrepareSync(machineName, port, rpath, rpathDir)
	Sync(machineName, port, path, pathpkg.Dir(path)) // initial sync

	if !*onetime {

		// catch ^C and restore vboxsf
		hitCounter := 0
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				if hitCounter > 0 {
					os.Exit(1)
				}

				hitCounter += 1
				fmt.Println()
				RestoreVBoxsf(machineName)
				os.Exit(0)
			}
		}()

		Watch(path, func(id uint64, path string, flags []string) {
			Sync(machineName, port, rpath, rpathDir)
		})
	} else {
		RestoreVBoxsf(machineName)
	}

}
