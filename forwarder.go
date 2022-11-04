package main

import "fmt"

type Forwarder func(from, to, addr string) error

type Resetter func() error

type Checker func() (string, error)

func forwarder(exec Executer) (Forwarder, Resetter, Checker) {
	const (
		portProxyCmd = "netsh interface portproxy add v4tov4 listenport=%s connectport=%s connectaddress=%s"
		firewallCmd  = "netsh advfirewall firewall add rule name=%s dir=in action=allow protocol=TCP localport=%s"
		resetterCmd  = "netsh interface portproxy reset"
		checkerCmd   = "netsh interface portproxy show v4tov4"
	)

	var (
		forwarder = func(from, to, addr string) error {
			_, err := exec(fmt.Sprintf(portProxyCmd, from, to, addr))
			if err != nil {
				return err
			}

			_, err = exec(fmt.Sprintf(firewallCmd, from, to))
			if err != nil {
				return err
			}

			return nil
		}

		resetter = func() error {
			_, err := exec(resetterCmd)
			if err != nil {
				return err
			}

			return nil
		}

		checker = func() (string, error) {
			cmd, err := exec(checkerCmd)
			if err != nil {
				return "", err
			}

			return cmd, nil
		}
	)

	return forwarder, resetter, checker
}
