package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func Provision(machineName string, verbose bool) {
	c := []string{
		// install and run rsync daemon
		`tce-load -wi rsync attr acl`,

		// disable boot2dockers builtin vboxfs
		// TODO bad idea, because you then can't use vboxfs anymore
		// `sudo umount /Users || /bin/true`,
	}

	for _, v := range c {
		out, err := RunSSHCommand(machineName, v, verbose)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("%s\n", out)
			os.Exit(1)
		}
	}
}

func RunSSHCommand(machineName, command string, verbose bool) (out []byte, err error) {
	if verbose {
		fmt.Println(`docker-machine ssh ` + machineName + ` '` + command + `'`)
	}
	return exec.Command("sh", "-c", `docker-machine ssh `+machineName+` '`+command+`'`).CombinedOutput()
}

func GetSSHPort(machineName string) (port uint, err error) {
	out, err := exec.Command("sh", "-c", `docker-machine inspect `+machineName).CombinedOutput()
	if err != nil {
		return 0, err
	}

	var v struct {
		Driver struct {
			SSHPort uint
		}
	}

	if err := json.Unmarshal(out, &v); err != nil {
		return 0, err
	}
	return v.Driver.SSHPort, nil
}
