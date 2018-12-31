package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"time"

	"github.com/dueckminor/mypi-api/go/config"
)

func Setup() {
	mypiRoot := config.GetRoot()
	mypiCfg := config.GetConfig()

	fmt.Println("mypiRoot:", mypiRoot)
	fmt.Println("mypiCfg:", mypiCfg)

	hostname := mypiCfg.GetString("config", "hostname")
	if len(hostname) == 0 {
		panic("hostname not configured")
	}

	ca, err := LoadCA(mypiRoot + "/config/pki/root_ca")
	if err != nil {
		ca, err = CreateRootCA("CN=MyPi-ROOT-CA")
		if err != nil {
			panic(err)
		}
		ca.Save(mypiRoot + "/config/pki/root_ca")
	}

	id, err := LoadIdentity(mypiRoot + "/config/pki/tls")
	if err != nil {
		template := &x509.Certificate{
			Subject:     pkix.Name{CommonName: hostname},
			DNSNames:    []string{hostname, "localhost"},
			IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().AddDate(20, 0, 0),
		}
		id = &Identity{}
		id.CreateKeyPair()
		id.certificate, err = ca.IssueCertificate(id, template)
		if err != nil {
			panic(err)
		}
		id.Save(mypiRoot + "/config/pki/tls")
	}

}
