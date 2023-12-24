package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"time"

	"github.com/dueckminor/mypi-tools/go/pki"
	"github.com/dueckminor/mypi-tools/go/util"
)

func updateCertificate(filename string, dNSNames []string) error {
	ca, err := pki.LoadCA("/opt/mypi/config/pki/root_ca")
	if err != nil {
		return err
	}

	template := &x509.Certificate{
		Subject:   pkix.Name{CommonName: dNSNames[0]},
		DNSNames:  dNSNames,
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),
	}

	id, err := pki.LoadIdentity(filename)
	if err == nil {
		fmt.Println(filename, id)
		if util.StringsContainsAll(id.X509Cert().DNSNames, dNSNames) {
			return nil
		}
	}

	id, err = pki.CreateIdentity()
	if err != nil {
		return err
	}
	_, err = ca.IssueCertificate(id, template)
	if err != nil {
		panic(err)
	}
	err = id.Save(filename)
	if err != nil {
		panic(err)
	}

	return nil
}

func main() {
	flag.Parse()
	var err error
	args := flag.Args()
	if args[0] == "pki" && args[1] == "update-certificate" {
		err = updateCertificate(args[2], args[3:])
	}

	if err != nil {
		panic(err)
	}
}
