package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SSHCredentials struct {
	IPAddress  string
	SSHPort    uint
	SSHUser    string
	SSHKeyPath string
}

func needsProvisioning(machineName string, verbose bool) bool {
	checkCommands := []string{
		`which rsync attr`,
	}

	for _, command := range checkCommands {
		if _, err := RunSSHCommand(machineName, command, verbose); err != nil {
			if verbose {
				fmt.Println("Provisioning the docker-machine")
			}
			return true
		}
	}

	return false
}

func Provision(machineName string, verbose bool) {
	if needsProvisioning(machineName, verbose) {
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
}

func RunSSHCommand(machineName, command string, verbose bool) (out []byte, err error) {
	if verbose {
		fmt.Println(`docker-machine ssh ` + machineName + ` '` + command + `'`)
	}
	return Exec("docker-machine", "ssh", machineName, command).CombinedOutput()
}

func GetSSHCredentials(machineName string) (creds SSHCredentials, err error) {
	out, err := Exec("docker-machine", "inspect", "--format='{{json .Driver}}'", machineName).CombinedOutput()
	if err != nil {
		return SSHCredentials{}, err
	}

	return CredentialsFromMachineJSON(out)
}

func CredentialsFromMachineJSON(jsonData []byte) (creds SSHCredentials, err error) {
	var v SSHCredentials
	if err := json.Unmarshal(jsonData, &v); err != nil {
		return SSHCredentials{}, err
	}

	return v, nil
}
