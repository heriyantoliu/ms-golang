package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"strings"
)

func ResolveIpFromHostsFile() (string, error) {
	data, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		logrus.Errorf("Problem reading /etc/hosts: %v", err.Error())
		return "", fmt.Errorf("Problem reading /etc/hosts: "+ err.Error())
	} else {
		lines := strings.Split(string(data), "\n")

		line := lines[len(lines) - 1]
		if len(line) < 2 {
			line = lines[len(lines)-2]
		}

		parts := strings.Split(line, "\t")
		return parts[0], nil
	}
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

func GetIPWithPrefix(prefix string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), prefix) {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}
