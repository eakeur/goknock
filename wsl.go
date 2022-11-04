package main

import (
	"errors"
	"strings"
)

type WSLIPGetter func() (string, error)

func wslIPGetter(exec Executer) WSLIPGetter {
	return func() (string, error) {
		result, err := exec("wsl ifconfig eth0")
		if err != nil {
			return "", err
		}

		wslIP, err := inetIP(result)
		if err != nil {
			return "", err
		}

		return wslIP, nil
	}
}

func inetIP(ifconfig string) (string, error) {
	var splitLine []string
	for _, ln := range strings.Split(ifconfig, "\n") {
		ln = strings.TrimSpace(ln)
		if !strings.HasPrefix(ln, "inet ") {
			continue
		}

		splitLine = strings.Split(ln, " ")
	}

	if len(splitLine) < 2 {
		return "", errors.New("failed fetching wsl ip from line. it is possible for you to provide it so this step won't be necessary")
	}

	return splitLine[1], nil
}
