package main

// docker-rsync start synack-dev

// switch to https://github.com/go-fsnotify/fsnotify once FSEvents is available

import (
	"fmt"
	"os"
)

func main() {
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
	PrepareSync(machineName, port, path, path)
	Sync(machineName, port, path, path)
	Watch(path, func(id uint64, path string, flags []string) {
		Sync(machineName, port, path, path)
	})
}
