package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	acOnlineFilePath         = "/sys/class/power_supply/ACAD/online"
	batteryModeCommand       = "/usr/bin/ryzenadj --power-saving && /usr/bin/cpupower frequency-set -g conservative"
	acModeCommand            = "/usr/bin/ryzenadj --max-performance && /usr/bin/cpupower frequency-set -g performance"
	pollingIntervalInSeconds = 10
)

func runCommandAsRoot(command string, sudoPassword string) {
	cmd := exec.Command("sudo", "-S", "bash", "-c", command)
	cmd.Stdin = bytes.NewBufferString(sudoPassword + "\n")
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to execute command '%s' as root: %v", command, err)
	} else {
		log.Printf("Command executed successfully: %s", command)
	}
}

func main() {
	battery := flag.Bool("battery", false, "force battery mode")
	plugged := flag.Bool("plugged", false, "force plugged mode")

	flag.Parse()
	fmt.Print("Enter sudo password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}

	sudoPassword := string(bytePassword)

	// State of the AC power, true means it's on, false means it's off
	var acState bool
	// Indicates if the acState was set before, used to determine when the state changes
	var acStateSet bool

	for {
		if *battery {
			log.Println("Forcing battery mode")
			runCommandAsRoot(batteryModeCommand, sudoPassword)
			*battery = false
		} else if *plugged {
			log.Println("Forcing plugged mode")
			runCommandAsRoot(acModeCommand, sudoPassword)
			*plugged = false
		} else {
			acOnline, err := os.ReadFile(acOnlineFilePath)
			if err != nil {
				log.Printf("Failed to read file '%s': %v", acOnlineFilePath, err)
			} else {
				newAcState := strings.TrimSpace(string(acOnline)) != "0"
				if !acStateSet || newAcState != acState {
					acState = newAcState
					acStateSet = true

					if acState {
						log.Println("On AC detected")
						runCommandAsRoot(acModeCommand, sudoPassword)
					} else {
						log.Println("On battery detected")
						runCommandAsRoot(batteryModeCommand, sudoPassword)
					}
				}
			}
		}

		time.Sleep(pollingIntervalInSeconds * time.Second)
	}
}
