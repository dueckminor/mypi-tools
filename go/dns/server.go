package dns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Server interface {
}

type server struct {
	mux *dns.ServeMux
	udp *dns.Server

	domains []string
	address string

	challenges map[string]string
}

func NewServer() (s Server, err error) {
	server := &server{}

	server.mux = dns.NewServeMux()
	server.mux.HandleFunc(".", server.dnsHandleFunc)

	server.udp = &dns.Server{
		Addr:    ":10053",
		Net:     "udp",
		Handler: server.mux,
	}

	server.address = "127.0.0.1"
	server.domains = []string{
		"rh94.dueckminor.de",
	}

	server.challenges = make(map[string]string)
	server.challenges["_acme-challenge.rh94.dueckminor.de."] = "foo"

	go func() {
		err := server.udp.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	return server, nil
}

func (s *server) Close() error {
	return s.udp.Shutdown()
}

func (s *server) matchHost(name string) bool {
	for _, domain := range s.domains {
		if strings.HasSuffix(name, "."+domain+".") {
			return true
		}
	}
	return false
}

func (s *server) dnsHandleFunc(w dns.ResponseWriter, r *dns.Msg) {

	m := new(dns.Msg)
	m.SetReply(r)

	ip := net.ParseIP(s.address)

	// t := &dns.TXT{
	// 	Hdr: dns.RR_Header{Name: "rh94.dueckminor.de", Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
	// 	Txt: []string{"something"},
	// }

	switch r.Question[0].Qtype {
	case dns.TypeTXT:
		if resp, ok := s.challenges[r.Question[0].Name]; ok {
			m.Answer = append(m.Answer, &dns.TXT{
				Hdr: dns.RR_Header{
					Name:   r.Question[0].Name,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				Txt: []string{resp},
			})
		}
	case dns.TypeA:
		fmt.Println("A", r.Question[0].Name)
		if s.matchHost(r.Question[0].Name) {
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   r.Question[0].Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: ip,
			})
		}
	case dns.TypeAAAA:
		fmt.Println("AAAA", r.Question[0].Name)
	default:
		fmt.Println("?", r.Question[0].Name)
	}

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			m.SetTsig(r.Extra[len(r.Extra)-1].(*dns.TSIG).Hdr.Name, dns.HmacSHA256, 300, time.Now().Unix())
		} else {
			println("Status", w.TsigStatus().Error())
		}
	}
	fmt.Printf("<<<<<<<<<<\n%v>>>>>>>>>>\n", m.String())

	err := w.WriteMsg(m)
	if err != nil {
		fmt.Println("failed to write DNS respones:", err)
	}
}
