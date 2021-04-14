package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dueckminor/mypi-tools/go/config"
)

type PkiGenerator struct {
	PkiDir string
	CA     CA
}

func (p *PkiGenerator) getFileName(name, certType string) string {
	name = strings.ReplaceAll(name, ".", "_")
	return path.Join(p.PkiDir, name+"_"+certType)
}

func (p *PkiGenerator) loadCa(name string) (ca CA, err error) {
	name = p.getFileName(name, "ca")
	return LoadCA(name)
}

func (p *PkiGenerator) GenerateRoot() (err error) {
	if p.CA != nil {
		return nil
	}
	p.CA, err = p.loadCa("root")
	if err == nil {
		return nil
	}

	err = os.MkdirAll(p.PkiDir, os.ModePerm)
	if err != nil {
		return err
	}

	p.CA, err = CreateRootCA("MyPi-ROOT-CA")
	if err != nil {
		return err
	}
	err = p.CA.Save(p.getFileName("root", "ca"))
	if err != nil {
		return err
	}
	return nil
}

func (p *PkiGenerator) GenerateServerCert(name string, dns ...string) (err error) {
	p.GenerateRoot()

	id, err := CreateIdentity()
	if err != nil {
		return err
	}

	template := &x509.Certificate{
		Subject:   pkix.Name{CommonName: name},
		DNSNames:  dns,
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),
	}

	p.CA.IssueCertificate(id, template)

	return id.Save(p.getFileName(name, "tls"))
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

	_, err = LoadIdentity(mypiRoot + "/config/pki/tls")
	if err != nil {
		template := &x509.Certificate{
			Subject:     pkix.Name{CommonName: hostname},
			DNSNames:    dNSNames,
			IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().AddDate(20, 0, 0),
		}
		id := &IdentityImpl{}
		id.CreateKeyPair()
		id.certificate, err = p.CA.IssueCertificate(id, template)
		if err != nil {
			panic(err)
		}
		id.Save(mypiRoot + "/config/pki/tls")
	}

}
