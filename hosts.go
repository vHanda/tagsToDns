package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const hostFilePath = "/etc/hosts"

// HostsFile is a simple structure for maniputalating the etc hosts file
type HostsFile struct{}

// Add an ip to the hosts file
func (h HostsFile) Add(ip string, hostname string) error {
	file, err := os.OpenFile(hostFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	text := ip + " " + hostname + "\n"
	_, err = file.WriteString(text)
	return err
}

// Remove an IP from the hosts file
func (h HostsFile) Remove(ip string) error {
	contents, err := ioutil.ReadFile(hostFilePath)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(hostFilePath, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range strings.Split(string(contents), "\n") {
		if len(line) == 0 {
			fmt.Fprintln(file)
			continue
		}

		if strings.HasPrefix(line, "#") || !strings.HasPrefix(line, ip) {
			fmt.Fprintln(file, line)
		}
	}

	return nil
}
