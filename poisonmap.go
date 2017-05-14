package dnspass

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

const dataURL = "https://raw.githubusercontent.com/freman/dnspass/master/data/poison.txt"

type poisonMap struct {
	sync.RWMutex
	list map[string]struct{}
}

func (m *poisonMap) UnmarshalTOML(decode func(interface{}) error) error {
	list := []string{}
	if err := decode(&list); err != nil {
		return err
	}

	m.list = make(map[string]struct{})
	for _, group := range list {
		m.list[group] = emptyStruct
	}

	return nil
}

func (m *poisonMap) pullUpdate(log logrus.FieldLogger) {
	log.Debug("Downloading poison host list")

	resp, err := http.Get(dataURL)
	if err != nil {
		log.WithError(err).Warn("Unable to download list")
		return
	}
	defer resp.Body.Close()

	n := make(map[string]struct{})
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		parts := strings.SplitN(strings.TrimSpace(scanner.Text()), "#", 2)
		sip := strings.TrimSpace(parts[0])

		if len(sip) == 0 {
			continue
		}

		ip := net.ParseIP(sip)
		if ip.String() != sip {
			continue
		}

		n[sip] = emptyStruct
	}

	if l := len(n); l > 0 {
		log.Debugf("Done, I know now %d poison hosts", l)
		m.Lock()
		defer m.Unlock()
		m.list = n
	}
}

func (m *poisonMap) isBad(group ...string) bool {
	m.RLock()
	defer m.RUnlock()
	for _, g := range group {
		if _, ok := m.list[g]; ok {
			return true
		}
	}
	return false
}
