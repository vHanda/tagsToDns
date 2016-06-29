package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	dnsDomain = ".service.consul"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		action := os.Getenv("SERF_EVENT")
		switch action {
		case "member-join":
		case "member-update":
			line := scanner.Text()
			ip, host, err := fetchIPAndHost(line)

			if err != nil {
				continue
			}

			hostsFile := HostsFile{}
			hostsFile.Remove(ip)
			hostsFile.Add(ip, host)
			break

		case "member-failed":
		case "member-leave":
		case "member-reap":
			line := scanner.Text()
			ip, _, err := fetchIPAndHost(line)

			if err != nil {
				continue
			}

			hostsFile := HostsFile{}
			hostsFile.Remove(ip)
			break

		default:
			panic("Invalid serf event: " + action)
		}

		fmt.Println("Processed:", scanner.Text())
	}
}

// TODO: Handle multiple hosts!
func fetchIPAndHost(line string) (string, string, error) {
	tagArray := strings.Split(line, "=")

	if len(tagArray) <= 1 {
		return "", "", errors.New("Invalid tag from serf")
	}

	data := strings.Split(tagArray[1], ":")
	if len(data) <= 2 {
		return "", "", errors.New("Invalid tag from serf")
	}

	ip := data[1]
	host := data[0]

	return ip, host, nil
}
