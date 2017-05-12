package main

import (
	"time"

	"github.com/urfave/cli"
)

const (
	svcName         = "dnspass"
	svcDesc         = "Bypass the Australian site blocking initiative."
	gracefulTimeout = 30 * time.Second
)

var (
	defaultFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-filename",
			Usage: "Config filename",
			Value: "config.toml",
		},
		cli.BoolFlag{
			Name:  "log-debug",
			Usage: "Debug level logging",
		},
	}
)
