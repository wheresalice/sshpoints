package main

import (
	"github.com/wheresalice/sshpoints/servers"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	(&cli.App{
		Commands: []*cli.Command{
			{
				Name:    "ssh",
				Aliases: []string{"s"},
				Usage:   "run the SSH server",
				Action: func(cCtx *cli.Context) error {
					servers.SSH(cCtx.String("redis"))
					return nil
				},
			},
			{
				Name:    "http",
				Aliases: []string{"h"},
				Usage:   "run the HTTP server",
				Action: func(cCtx *cli.Context) error {
					servers.HTTP(cCtx.String("redis"))
					return nil
				},
			},
		},
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name: "redis",
        Usage: "redis connection string",
        Value: "localhost:6379",
      },
    },
	}).Run(os.Args)
}
