package setup

import (
	"os"
	"path"

	"github.com/dueckminor/mypi-tools/go/pki"
)

type Certificates struct {
	pkiDir pki.PKIDir
}

type Certificate struct {
	Label          string `json:"label"`
	Subject        string `json:"subject"`
	ValidNotAfter  string `json:"valid_not_after"`
	ValidNotBefore string `json:"valid_not_before"`
}

func NewCertificates() *Certificates {
	pkiDir, _ := pki.NewPKIDir(path.Join(os.Getenv("HOME"), ".mypi/pki"))
	return &Certificates{
		pkiDir: pkiDir,
	}
}

func (c *Certificates) Get() (certs []*Certificate, err error) {
	identities, err := c.pkiDir.GetIdentities()
	if err != nil {
		return nil, err
	}
	for _, identity := range identities {
		cert := &Certificate{}
		cert.Label = identity.Label()
		x509 := identity.X509Cert()
		cert.ValidNotBefore = x509.NotBefore.Format("2006-01-02")
		cert.ValidNotAfter = x509.NotAfter.Format("2006-01-02")
		cert.Subject = x509.Subject.String()
		certs = append(certs, cert)
	}
	return certs, nil
}

func (c *Certificates) CreatePKI() (certs []*Certificate, err error) {
	generator := c.pkiDir.GetGenerator()
	err = generator.GenerateRoot()
	if err != nil {
		return nil, err
	}
	err = generator.GenerateServerCert("localhost", "localhost")
	if err != nil {
		return nil, err
	}
	err = generator.GenerateServerCert("mypi.fritz.box", "mypi.fritz.box")
	if err != nil {
		return nil, err
	}
	return c.Get()
}
