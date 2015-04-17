package main

import (
	"fmt"
	"os"
	pathpkg "path"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println(Version)
		os.Exit(0)
	}

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("error: unable to get current directory:", err)
		os.Exit(1)
	}

	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Printf("usage: %v <machine-name>\n", os.Args[0])
		os.Exit(1)
	}

	machineName := os.Args[1]

	port, err := GetSSHPort(machineName)
	if err != nil {
		fmt.Println("error: unable to get port:", err)
		os.Exit(1)
	}

	Provision(machineName)

	rpath := path
	rpathDir := pathpkg.Dir(path)

	PrepareSync(machineName, port, rpath, rpathDir)
	Sync(machineName, port, path, pathpkg.Dir(path)) // initial sync
	Watch(path, func(id uint64, path string, flags []string) {
		Sync(machineName, port, rpath, rpathDir)
	})
}
