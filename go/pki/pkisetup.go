package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"path"
	"time"

	"github.com/dueckminor/mypi-tools/go/config"
)

type PkiGenerator struct {
	PkiDir string
	CA     *CA
}

func (p *PkiGenerator) GenerateRoot() (err error) {
	dir := path.Join(p.PkiDir, "root_ca")
	p.CA, err = LoadCA(dir)
	if err == nil {
		return nil
	}

	p.CA, err = CreateRootCA("MyPi-ROOT-CA")
	if err != nil {
		return err
	}
	err = p.CA.Save(dir)
	if err != nil {
		return err
	}
	return nil
}

func Setup() {
	mypiRoot := config.GetRoot()
	mypiCfg := config.GetConfig()

	fmt.Println("mypiRoot:", mypiRoot)
	fmt.Println("mypiCfg:", mypiCfg)

	hostname := mypiCfg.GetString("config", "hostname")
	if len(hostname) == 0 {
		panic("hostname not configured")
	}

	p := &PkiGenerator{}
	p.PkiDir = path.Join(mypiRoot, "/config/pki")
	err := p.GenerateRoot()
	if err != nil {
		panic(err)
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
		id.certificate, err = p.CA.IssueCertificate(id, template)
		if err != nil {
			panic(err)
		}
		id.Save(mypiRoot + "/config/pki/tls")
	}

}
