package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var lastSyncError = ""

func Sync(via string, port uint, src, dst string, verbose bool) {
	args := []string{
		// "--verbose",
		// "--stats",
		"--recursive",
		"--links",
		"--times",
		"--inplace",
		"--itemize-changes",
		"--delete",
		"--force",
		"--executability",
		"--compress",
	}

	ripath := getRsyncIgnorePath()
	if ripath != "" {
		args = append(args, `--exclude-from='`+ripath+`'`)
	}

	if strings.HasPrefix(via, "rsync://") {
		args = append(args, filepath.Join(src)+"/.")
		args = append(args, via)
	} else {
		machineName := via
		u, err := user.Current()
		if err != nil {
			panic(fmt.Sprintf("Unable to load current user's profile: %s", err))
		}

		sshKeyFile := filepath.Join(u.HomeDir, "/.docker/machine/machines/", machineName, "id_rsa")
		sshArg := fmt.Sprintf(`-e "ssh -o StrictHostKeyChecking=no -i %s -p %v"`, sshKeyFile, port)

		args = append(args, sshArg)
		args = append(args, "--rsync-path='sudo rsync'")
		args = append(args, src, "docker@localhost:"+dst)
	}

	cmd := Exec("rsync", args...)

	if verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// don't show duplicate errors
		if lastSyncError != err.Error() {
			fmt.Printf("error: %v\n", err)
		}
		lastSyncError = err.Error()
	}
}

func getRsyncIgnorePath() string {
	if _, err := os.Stat(".rsyncignore"); err == nil {
		abs, err := filepath.Abs(".rsyncignore")
		if err == nil {
			return abs
		}
	}
	return ""
}
