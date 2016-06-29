package main

import (
	"bufio"
	"errors"
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
			e, err := parse(scanner.Text())
			if err != nil {
				continue
			}

			hostsFile := HostsFile{}
			hostsFile.Remove(e.ip)
			hostsFile.Add(e.ip, addDNSDomain(e.hosts))
			break

		case "member-failed":
		case "member-leave":
		case "member-reap":
			e, err := parse(scanner.Text())
			if err != nil {
				continue
			}

			hostsFile := HostsFile{}
			hostsFile.Remove(e.ip)
			break

		default:
			panic("Invalid serf event: " + action)
		}
	}
}

// SerfEvent groups together the information given by serf
type SerfEvent struct {
	ip    string
	hosts []string
}

func parse(line string) (SerfEvent, error) {
	tagArray := strings.Split(line, "=")

	if len(tagArray) <= 1 {
		return SerfEvent{}, errors.New("Invalid tag from serf")
	}

	data := strings.Split(tagArray[1], ":")
	if len(data) <= 2 {
		return SerfEvent{}, errors.New("Invalid tag from serf")
	}

	return SerfEvent{data[1], strings.Split(data[0], ",")}, nil
}

func addDNSDomain(hosts []string) []string {
	const dnsDomain = ".service.consul"

	domains := make([]string, len(hosts))
	for i, v := range hosts {
		domains[i] = v + dnsDomain
	}
	return domains
}
