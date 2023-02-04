package network

import (
	"bufio"
	"net"
	"os/exec"
	"strings"
)

func init() {
	registerCbRouterInternalIP(getRouterInternalIPDarwin)
}

func getRouterInternalIPDarwin() (routerIP *net.IPAddr, err error) {
	cmd := exec.Command("route", "-n", "get", "default")

	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {
		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimLeft(line, " \t")
			if strings.HasPrefix(line, "router:") {
				line = strings.Trim(line[8:], " \t")
				routerIP, err = net.ResolveIPAddr("", line)
				break
			}
			if strings.HasPrefix(line, "gateway:") {
				line = strings.Trim(line[9:], " \t")
				routerIP, err = net.ResolveIPAddr("", line)
				break
			}
		}
		// We're all done, unblock the channel
		done <- struct{}{}
	}()

	// Start the command and check for errors
	err = cmd.Start()

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	err = cmd.Wait()

	return routerIP, err
}
