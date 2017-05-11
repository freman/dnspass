// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"golang.org/x/sys/windows/svc/debug"
)

// LogHook to send logs via windows log.
type WinLogHook struct {
	Upstream debug.Log
}

func NewWinLogHook(logger debug.Log) *WinLogHook {
	return &WinLogHook{logger}
}

func (hook *WinLogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		hook.Upstream.Error(3, line)
		os.Exit(1)
		return nil
	case logrus.FatalLevel:
		return hook.Upstream.Error(2, line)
		os.Exit(1)
		return nil
	case logrus.ErrorLevel:
		return hook.Upstream.Error(1, line)
	case logrus.WarnLevel:
		return hook.Upstream.Warning(1, line)
	case logrus.InfoLevel:
		return hook.Upstream.Info(2, line)
	case logrus.DebugLevel:
		return hook.Upstream.Info(1, line)
	default:
		return nil
	}
}

func (hook *WinLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
