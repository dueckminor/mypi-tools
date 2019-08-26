package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"time"

	"github.com/dueckminor/mypi-tools/go/config"
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
		ca, err = CreateRootCA("MyPi-ROOT-CA")
		if err != nil {
			panic(err)
		}
		err = ca.Save(mypiRoot + "/config/pki/root_ca")
		if err != nil {
			panic(err)
		}
	}

	dNSNames := []string{hostname}
	for _, hostname := range mypiCfg.GetArray("config", "hostnames") {
		dNSNames = append(dNSNames, hostname.GetString())
	}

	id, err := LoadIdentity(mypiRoot + "/config/pki/tls")
	if err != nil {
		template := &x509.Certificate{
			Subject:     pkix.Name{CommonName: hostname},
			DNSNames:    dNSNames,
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
