package main

import (
	"fmt"
	"gopkg.in/mattes/go-expand-tilde.v1"
	"os/exec"
)

func PrepareSync(machineName string, port uint, src, dst string) {
	RunSSHCommand(machineName, "sudo mkdir -p "+dst)
}

func Sync(machineName string, port uint, src, dst string) {
	homePath, _ := tilde.Expand("~/")
	out, err := exec.Command("/bin/sh", "-c", `rsync -rzv --force --exclude '.git/*' -e "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet -i `+homePath+`.docker/machine/machines/`+machineName+`/id_rsa -p `+fmt.Sprintf("%v", port)+`" --rsync-path="sudo rsync" `+src+` docker@localhost:`+dst).CombinedOutput()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("%s\n", out)
}
