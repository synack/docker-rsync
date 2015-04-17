package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func PrepareSync(machineName string, port uint, src, dst string) {
	RunSSHCommand(machineName, "sudo mkdir -p "+dst)
}

func Sync(machineName string, port uint, src, dst string) {
	homePath := os.Getenv("HOME")
	ripath := getRsyncIgnorePath()

	args := []string{
		"--recursive",
		"--links",
		"--times",
		"--inplace",
		// "--verbose",
		// "--stats",
		"--itemize-changes",
		"--delete",
		"--force",
		"--executability",
		"--compress",
		"--force",
		fmt.Sprintf(`-e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet -i "%s" -p %v'`, filepath.Join(homePath, "/.docker/machine/machines", machineName, "id_rsa"), port),
		"--rsync-path='sudo rsync'",
	}
	if ripath != "" {
		args = append(args, `--exclude-from='`+ripath+`'`)
	}
	args = append(args, src, "docker@localhost:"+dst)

	command := "rsync " + strings.Join(args, " ")

	// fmt.Println("/bin/sh", "-c", command)
	cmd := exec.Command("/bin/sh", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
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
