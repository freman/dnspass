// +build linux darwin freebsd

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/freman/dnspass"
	"github.com/urfave/cli"
)

func runService(c *cli.Context) error {
	log := logrus.New()
	if c.Bool("log-debug") {
		log.Level = logrus.DebugLevel
	}

	configFilename := c.String("config-filename")
	if configFilename == "config.toml" {
		exePath, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Fatalf("Failed to find the dnspass executable: %s", err)
		}
		configFilename = filepath.Join(filepath.Dir(exePath), configFilename)
	}

	server := &dnspass.Server{
		ConfigFile: configFilename,
		Log:        log,
	}

	var gracefulStop = make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	err := server.Run()
	if err != nil {
		log.Fatalf("Service failed to run: %s", err)
	}

	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v", sig)

	time.AfterFunc(gracefulTimeout, func() {
		log.Error("Failed to stop quickly; stopping forcefully")
		os.Exit(1)
	})

	server.Shutdown()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = svcName
	app.Usage = svcName + " - " + svcDesc
	app.Version = dnspass.Version
	app.Flags = defaultFlags
	app.Action = runService
	app.Run(os.Args)
}
