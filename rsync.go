package main

import (
	"fmt"
	"os"
	"os/exec"
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
		homePath := os.Getenv("HOME")
		args = append(args, fmt.Sprintf(`-e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet -i "%s" -p %v'`, filepath.Join(homePath, "/.docker/machine/machines", machineName, "id_rsa"), port))
		args = append(args, "--rsync-path='sudo rsync'")
		args = append(args, src, "docker@localhost:"+dst)
	}

	command := "rsync " + strings.Join(args, " ")

	// fmt.Println("/bin/sh", "-c", command)
	cmd := exec.Command("/bin/sh", "-c", command)

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
