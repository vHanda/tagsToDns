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

			hostsFile := NewHostsFile()
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

			hostsFile := NewHostsFile()
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
	id    string
	hosts []string
}

func parse(line string) (SerfEvent, error) {
	var event SerfEvent

	tags := strings.Split(line, ",")
	for _, tag := range tags {
		tagArray := strings.Split(tag, "=")

		if len(tagArray) <= 1 {
			return SerfEvent{}, errors.New("Invalid tag from serf")
		}

		data := strings.Split(tagArray[1], ":")
		if len(data) <= 2 {
			return SerfEvent{}, errors.New("Invalid tag from serf")
		}

		// The IP and id are being overwritten when there is more than 1
		// but in practice they are the same, so it doesn't matter
		event.ip = data[1]
		event.id = tagArray[0]
		event.hosts = append(event.hosts, data[0])
	}

	return event, nil
}

func addDNSDomain(hosts []string) []string {
	const dnsDomain = ".service.consul"

	domains := make([]string, len(hosts))
	for i, v := range hosts {
		domains[i] = v + dnsDomain
	}
	return domains
}
