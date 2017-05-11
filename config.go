package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/naoina/toml"
)

type hostMap map[string]struct{}

var emptyStruct = struct{}{}

var config = struct {
	Listen   string
	Trust    []string
	Untrust  []string
	BadHosts hostMap
}{
	Listen: "localhost:53",
	Untrust: []string{
		"203.12.160.35:53",
		"203.12.160.36:53",
	},
	Trust: []string{
		"8.8.8.8:53",
		"8.8.4.4:53",
	},
	BadHosts: hostMap{
		"202.136.99.184": emptyStruct,
		"202.136.99.185": emptyStruct,
		"101.167.166.53": emptyStruct,
		"54.79.39.115":   emptyStruct,
	},
}

func (m *hostMap) UnmarshalTOML(decode func(interface{}) error) error {
	list := []string{}
	if err := decode(&list); err != nil {
		return err
	}

	*m = make(hostMap)
	for _, group := range list {
		(*m)[group] = emptyStruct
	}

	return nil
}

func (m *hostMap) isBad(group ...string) bool {
	for _, g := range group {
		if _, ok := (*m)[g]; ok {
			return true
		}
	}
	return false
}

func loadConfig(log logrus.FieldLogger, fn string) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		log.Fatal(err)
	}
}
