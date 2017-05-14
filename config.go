package dnspass

import (
	"errors"
	"os"
	"strings"

	"github.com/naoina/toml"
)

var emptyStruct = struct{}{}

type config struct {
	AutoUpdatePoisonHosts bool
	Listen                string
	Trust                 []string
	Untrust               []string
	PoisonHosts           poisonMap
}

func (s *Server) loadConfig() error {
	f, err := os.Open(s.ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := toml.NewDecoder(f).Decode(&s.config); err != nil {
		return err
	}

	if len(s.config.Trust) < 1 {
		return errors.New("Need at least one trusted dns server")
	}
	if len(s.config.Untrust) < 1 {
		return errors.New("Need at least one untrusted dns server")
	}

	// Set the peers to be port 53 if not already
	for i := range s.config.Trust {
		if !strings.Contains(s.config.Trust[i], ":") {
			s.config.Trust[i] += ":53"
		}
	}
	for i := range s.config.Untrust {
		if !strings.Contains(s.config.Trust[i], ":") {
			s.config.Trust[i] += ":53"
		}
	}

	return nil
}
