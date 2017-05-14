package dnspass

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	lru "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
)

var (
	Version = "Undefined"
	Commit  = "Undefined"
)

// Server is the all the things a good server need
type Server struct {
	Log        logrus.FieldLogger
	ConfigFile string
	config     config
	cache      *lru.Cache
	rand       *rand.Rand
	servers    []*dns.Server
	stopping   bool
}

// Run starts two listeners - one for tcp and one for udp
func (s *Server) Run() error {
	err := s.loadConfig()
	if err != nil {
		return fmt.Errorf("Unable to load configuration: %s", err)
	}

	s.cache, err = lru.New(256)
	if err != nil {
		return fmt.Errorf("Unable to construct cache: %s", err)
	}

	if s.config.AutoUpdatePoisonHosts {
		go s.config.PoisonHosts.pullUpdate(s.Log)
		go func() {
			for {
				time.Sleep(24 * time.Hour)
				s.config.PoisonHosts.pullUpdate(s.Log)
			}
		}()
	}

	s.setupRand()

	s.servers = []*dns.Server{}

	mux := dns.NewServeMux()
	mux.Handle(".", s)

	for _, proto := range []string{"udp", "tcp"} {
		go func(proto string) {
			l := s.Log.WithFields(logrus.Fields{
				"listen": s.config.Listen,
				"proto":  proto,
			})

			server := &dns.Server{
				Addr:    s.config.Listen,
				Net:     proto,
				Handler: mux,
			}

			s.servers = append(s.servers, server)

			if err := server.ListenAndServe(); err != nil {
				l.WithError(err).Fatal("Unable to listen")
			}

			if !s.stopping {
				l.Fatal("Server exited expectantly")
			}
		}(proto)
	}

	return nil
}

// Shutdown the listeners
func (s *Server) Shutdown() {
	s.stopping = true
	for _, server := range s.servers {
		server.Shutdown()
	}
}

// ServeDNS implements the dns.Handler interface
func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log := s.Log.WithField("question", r.Question)

	proto := "tcp"
	if _, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		proto = "udp"
	}

	if len(r.Question) != 1 || r.Question[0].Qtype != dns.TypeA {
		log.Debug("Question too complicated, just forwarding")
		w.WriteMsg(s.forwardRequest(s.config.Untrust, proto, r))
		return
	}

	q := r.Question[0]

	// Is the request known to result in poisoned DNS response?
	if _, found := s.cache.Get(q); !found {
		in := s.forwardRequest(s.config.Untrust, proto, r)
		if in.Rcode == dns.RcodeNameError {
			w.WriteMsg(in)
			return
		}

		poisoned := false
		for _, ans := range in.Answer {
			if a, ok := ans.(*dns.A); ok && s.config.PoisonHosts.isBad(a.A.String()) {
				poisoned = true
				break
			}
		}

		// Not poisoned, good reply
		if !poisoned {
			w.WriteMsg(in)
			return
		}

		log.Warn("ISP Poisoned answer detected")
		s.cache.Add(q, emptyStruct)
	}

	log.Debug("Sending to trusted DNS servers")
	w.WriteMsg(s.forwardRequest(s.config.Trust, proto, r))
}

func (s *Server) setupRand() {
	s.rand = rand.New(rand.NewSource(time.Now().Unix()))
}

func (s *Server) forwardRequest(peers []string, proto string, r *dns.Msg) (in *dns.Msg) {
	var err error

	// Try each server of the peer servers in order until one works or they all fail
	for _, i := range s.rand.Perm(len(peers)) {
		c := new(dns.Client)
		c.Net = proto
		in, _, err = c.Exchange(r, peers[i])

		// One worked, yay
		if err == nil {
			return in
		}
	}

	s.Log.WithFields(logrus.Fields{
		"peers":     peers,
		"proto":     proto,
		"request":   r,
		"lastError": err,
	}).Warn("Failed to get answer from peers")

	in = new(dns.Msg)
	in.SetRcode(r, dns.RcodeServerFailure)
	return in
}
