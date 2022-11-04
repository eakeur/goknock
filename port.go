package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Ports map[string]string

func fetchPorts(pairs []string) (Ports, error) {
	if len(pairs) == 0 {
		return nil, errors.New("at least one port should be provided")
	}

	ports := Ports{}
	for _, arg := range pairs {
		pair := strings.Split(arg, ":")
		for _, p := range pair {
			_, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid port value: %s", p)
			}
		}

		if len(pair) < 2 {
			ports[pair[0]] = pair[0]
			continue
		}

		ports[pair[0]] = pair[1]
	}

	return ports, nil
}
