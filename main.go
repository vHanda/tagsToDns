package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
			removeFromHostsFile(ip)
			addToHostsFile(ip, host+".service.consul")
			break

		case "member-failed":
		case "member-leave":
		case "member-reap":
			line := scanner.Text()
			ip, _, err := fetchIPAndHost(line)

			if err != nil {
				continue
			}
			removeFromHostsFile(ip)
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

func addToHostsFile(ip, hostname string) {
	filePath := "/tmp/hosts"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	text := ip + " " + hostname
	if _, err = file.WriteString(text); err != nil {
		panic(err)
	}
}

// This is shitty. I'm sure it can be simplified far more!
func removeFromHostsFile(ip string) {
	filePath := "/tmp/hosts"

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
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
}
