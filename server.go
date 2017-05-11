package main

import (
	"math/rand"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	lru "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
)

var (
	cache *lru.Cache
)

func init() {
	rand.Seed(time.Now().Unix())
}

func doServer(configPath string, log logrus.FieldLogger) {
	loadConfig(log, configPath)
	cache, _ = lru.New(256)
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		log := log.WithField("question", r.Question)

		proto := "tcp"
		if _, ok := w.RemoteAddr().(*net.UDPAddr); ok {
			proto = "udp"
		}

		if len(r.Question) != 1 || r.Question[0].Qtype != dns.TypeA {
			log.Debug("Question too complicated, just forwarding")
			w.WriteMsg(forwardRequest(log, config.Untrust, proto, r))
			return
		}

		q := r.Question[0]

		if _, found := cache.Get(q); !found {
			in := forwardRequest(log, config.Untrust, proto, r)
			if in.Rcode == dns.RcodeNameError {
				w.WriteMsg(in)
				return
			}

			bad := false
			for _, ans := range in.Answer {
				if a, ok := ans.(*dns.A); ok && config.BadHosts.isBad(a.A.String()) {
					bad = true
					break
				}
			}

			if !bad {
				w.WriteMsg(in)
				return
			}

			log.Debug("Eww, ISP jiggered it")
			cache.Add(q, emptyStruct)
		}

		log.Debug("Sending to trusted dns servers")
		w.WriteMsg(forwardRequest(log, config.Trust, proto, r))
	})

	for _, proto := range []string{"udp", "tcp"} {
		go func(proto string) {
			l := log.WithFields(logrus.Fields{
				"listen": config.Listen,
				"proto":  proto,
			})
			server := &dns.Server{Addr: config.Listen, Net: proto}
			if err := server.ListenAndServe(); err != nil {
				l.WithError(err).Fatal("Unable to listen")
			}

			l.Fatal("Server exited expectantly")
		}(proto)
	}
}

func forwardRequest(log logrus.FieldLogger, upstream []string, proto string, r *dns.Msg) (in *dns.Msg) {
	var err error
	for _, i := range rand.Perm(len(upstream)) {
		c := new(dns.Client)
		c.Net = proto
		in, _, err = c.Exchange(r, upstream[i])
		if err == nil {
			return
		}
	}

	log.WithFields(logrus.Fields{
		"upstream": upstream,
		"proto":    proto,
		"r":        r,
		"lastErr":  err,
	}).Debug("Failed to get answer from upstream")

	in = new(dns.Msg)
	in.SetRcode(r, dns.RcodeServerFailure)
	return
}
