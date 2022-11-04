package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	powershell := executer()
	forward, reset, check := forwarder(powershell)
	getLinuxIPV4 := wslIPGetter(powershell)

	app := &cli.App{
		Name:        "Goknock",
		Description: "Make port forwarding on your windows machine easier",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "reset",
				Aliases: []string{"r"},
				Usage:   "Defines if the actual port forwarding configuration should be reset.",
			},
			&cli.StringFlag{
				Name:    "ip",
				Aliases: []string{"i"},
				Usage:   "Provides the IP address to bind to.",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "reset",
				Aliases: []string{"r"},
				Usage:   "Removes the actual port forwarding configurations",
				Action: func(ctx *cli.Context) error {
					fmt.Println("Resetting port forwarding configuration")
					return reset()
				},
			},
			{
				Name:    "wsl",
				Aliases: []string{"r"},
				Usage:   "Gets the default WSL container's IP address using the linux ifconfig command.",
				Action: func(ctx *cli.Context) error {
					fmt.Println("Attempting to fetch a WSL container's IP address")
					ip, err := getLinuxIPV4()
					if err != nil {
						return err
					}

					fmt.Println(ip)

					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) (err error) {

			ports, err := fetchPorts(ctx.Args().Slice())
			if err != nil {
				return err
			}

			ip := ctx.String("ip")
			if ip == "" {
				fmt.Println("Fetching WSL container's IP address")
				if ip, err = getLinuxIPV4(); err != nil {
					return err
				}
			}

			if ctx.Bool("reset") {
				fmt.Println("Resetting port forwarding configuration")
				if err := reset(); err != nil {
					return err
				}
			}

			for from, to := range ports {
				fmt.Printf("Forwarding port %s to %s\n", from, to)
				if err := forward(from, to, ip); err != nil {
					fmt.Printf("Failed forwarding port %s to %s\n: %v", from, to, err)
				}
			}

			res, err := check()
			if err != nil {
				return err
			}

			fmt.Println(res)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
