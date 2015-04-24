package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func Provision(machineName string) {
	c := []string{
		// install and run rsync daemon
		`tce-load -wi rsync`,

		// disable boot2dockers builtin vboxsf
		`sudo umount /Users || /bin/true`,
	}

	for _, v := range c {
		out, err := RunSSHCommand(machineName, v)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("%s\n", out)
			os.Exit(1)
		}
	}
}

func RestoreVBoxsf(machineName string) {
	c := []string{
		// disable rsync share
		`sudo umount /Users || /bin/true`,

		// restore vboxsf
		`sudo /etc/rc.d/automount-shares || /bin/true`,
	}

	for _, v := range c {
		out, err := RunSSHCommand(machineName, v)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("%s\n", out)
			os.Exit(1)
		}
	}
}

func RunSSHCommand(machineName, command string) (out []byte, err error) {
	fmt.Println(`docker-machine ssh ` + machineName + ` '` + command + `'`)
	return exec.Command("/bin/sh", "-c", `docker-machine ssh `+machineName+` '`+command+`'`).CombinedOutput()
}

func GetSSHPort(machineName string) (port uint, err error) {
	out, err := exec.Command("/bin/sh", "-c", `docker-machine inspect `+machineName).CombinedOutput()
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
