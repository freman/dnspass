// +build darwin linux

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	return "", err
}

func run() {
	debug := flag.Bool("debug", false, "Debug log level")
	showVersion := flag.Bool("version", false, "Show version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("dnspass - %s (%s)\n", version, commit)
		fmt.Println("https://github.com/freman/dnspass")
		return
	}

	log := logrus.New()

	if *debug {
		log.Level = logrus.DebugLevel
	}

	doServer("config.toml", log)

	ch := make(chan bool, 1)
	<-ch
}
