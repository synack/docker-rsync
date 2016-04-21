// +build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func Exec(cmd string, args ...string) *exec.Cmd {
	shCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	a := []string{"-Command", shCmd}

	return exec.Command("powershell", a...)
}
