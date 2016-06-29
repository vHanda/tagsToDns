package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// HostsFile is a simple structure for maniputalating the etc hosts file
type HostsFile struct {
	filePath string
}

// NewHostsFile creates a new object
func NewHostsFile() HostsFile {
	var h HostsFile
	h.filePath = os.Getenv("DISCOVERY_HOSTS_FILE_PATH")
	if h.filePath == "" {
		h.filePath = "/hosts/hosts.serf"
	}

	return h
}

// Add an ip to the hosts file
func (h HostsFile) Add(ip string, hosts []string) error {
	file, err := os.OpenFile(h.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Println(ip, hosts)
	text := ip + " " + strings.Join(hosts, " ")
	fmt.Println(text + "FOO")
	_, err = file.WriteString(text)
	return err
}

// Remove an IP from the hosts file
func (h HostsFile) Remove(ip string) error {
	contents, _ := ioutil.ReadFile(h.filePath)

	file, err := os.OpenFile(h.filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
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
