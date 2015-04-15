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

		// add init script
		`echo "#/bin/sh" | sudo tee /var/lib/boot2docker/bootlocal.sh`,
		`echo "sudo umount /Users" | sudo tee -a /var/lib/boot2docker/bootlocal.sh`,
		`sudo chmod +x /var/lib/boot2docker/bootlocal.sh`,

		// make changes live without rebooting
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

func RunSSHCommand(machineName, command string) (out []byte, err error) {
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
