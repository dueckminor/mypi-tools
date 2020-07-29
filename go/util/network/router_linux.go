package network

import (
	"bufio"
	"encoding/hex"
	"net"
	"os"
	"strings"
)

func init() {
	registerCbRouterInternalIP(getRouterInternalIPLinux)
}

func getRouterInternalIPLinux() (routerIP *net.IPAddr, err error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	scanner.Scan()
	parts := strings.Fields(scanner.Text())
	if len(parts) < 3 || len(parts[2]) != 8 {
		return nil, nil
	}

	ip, err := hex.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}
	ip[0], ip[1], ip[2], ip[3] = ip[3], ip[2], ip[1], ip[0]
	return &net.IPAddr{IP: ip}, nil
}
