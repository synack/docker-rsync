package main

import "testing"

func TestPortFromMachineJSONOldStyle(t *testing.T) {
    const jsonData string = `
        {
            "Driver": {
                "SSHPort": 52435
            }
        }
    `
    out, _ := PortFromMachineJSON([]byte(jsonData))
    if out != 52435 {
        t.Errorf("Expecting 52435 and got %d", out)
    }
}

func TestPortFromMachineJSONNewStyle(t *testing.T) {
    const jsonData string = `
        {
            "Driver": {
                "Driver": {
                    "SSHPort": 52435
                }
            }
        }
    `
    out, _ := PortFromMachineJSON([]byte(jsonData))
    if out != 52435 {
        t.Errorf("Expecting 52435 and got %d", out)
    }
}
