package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"time"
)

// CA can be used to issue other certs
type Identity interface {
	Label() string
	X509Cert() *x509.Certificate
	PublicKey() *rsa.PublicKey
	PrivateKey() *rsa.PrivateKey
	Save(filename string) error
}
type CA interface {
	Identity
	IssueCertificate(pub interface{}, template *x509.Certificate) (*x509.Certificate, error)
}

type IdentityImpl struct {
	label       string
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	certificate *x509.Certificate
}

type CAImpl struct {
	IdentityImpl
}

func LoadIdentity(filename string) (id_ Identity, err error) {
	id := &IdentityImpl{}
	err = id.Load(filename)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func LoadCA(filename string) (ca_ CA, err error) {
	ca := &CAImpl{}
	err = ca.Load(filename)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (id *IdentityImpl) Label() string {
	return id.label
}
func (id *IdentityImpl) X509Cert() *x509.Certificate {
	return id.certificate
}
func (id *IdentityImpl) PublicKey() *rsa.PublicKey {
	return id.publicKey
}
func (id *IdentityImpl) PrivateKey() *rsa.PrivateKey {
	return id.privateKey
}

func (id *IdentityImpl) Load(filename string) error {
	id.label = filepath.Base(filename)
	pem, err := ioutil.ReadFile(filename + "_priv.pem")
	if err != nil {
		return err
	}
	id.privateKey, err = ParseRsaPrivateKeyFromPem(pem)
	if err != nil {
		return err
	}
	id.publicKey = &id.privateKey.PublicKey
	pem, err = ioutil.ReadFile(filename + "_cert.pem")
	if err != nil {
		return err
	}
	id.certificate, err = ParseCertificateFromPem(pem)
	if err != nil {
		return err
	}
	return nil
}

func (id *IdentityImpl) Save(filename string) error {
	if id.privateKey != nil {
		binary := x509.MarshalPKCS1PrivateKey(id.privateKey)
		err := ioutil.WriteFile(filename+"_priv.pem", pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: binary,
		}), 0600)
		if err != nil {
			return err
		}
	}
	if id.certificate != nil {
		err := ioutil.WriteFile(filename+"_cert.pem", pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: id.certificate.Raw,
		}), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func RsaPrivateKeyToPem(priv *rsa.PrivateKey) []byte {
	binary := x509.MarshalPKCS1PrivateKey(priv)
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: binary})
}

func RsaPublicKeyToPem(pub *rsa.PublicKey) []byte {
	binary, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		panic(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: binary})
}

func ParseRsaPrivateKeyFromPem(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ParseCertificateFromPem(privPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (id *IdentityImpl) CreateKeyPair() {
	id.privateKey, id.publicKey = GenerateRsaKeyPair()
}

func CreateIdentity() (id Identity, err error) {
	res := &IdentityImpl{}
	res.CreateKeyPair()
	return res, nil
}

// CreateRootCA creates a root ca
func CreateRootCA(cn string) (ca CA, err error) {
	res := &CAImpl{}
	res.CreateKeyPair()

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: cn,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(20, 0, 0),
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, res.publicKey, res.privateKey)
	if err != nil {
		return
	}

	res.certificate, err = x509.ParseCertificate(cert)
	if err != nil {
		return
	}

	return res, nil
}

func (ca *CAImpl) IssueCertificate(pub interface{}, template *x509.Certificate) (*x509.Certificate, error) {
	id, ok := pub.(*IdentityImpl)
	if ok {
		pub = id.publicKey
	}

	if nil == template.SerialNumber {
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
		if err != nil {
			return nil, err
		}
		template.SerialNumber = serialNumber
	}

	certBuffer, err := x509.CreateCertificate(rand.Reader, template, ca.certificate, pub, ca.privateKey)
	if err != nil {
		return nil, err
	}
	id.certificate, err = x509.ParseCertificate(certBuffer)
	return id.certificate, err
}

// // IssueCertificate issues a certificate
// func (ca *CA) IssueCertificate(pemCSR string) (certsPEM string, err error) {
// 	csr, err := PEMToCSR(pemCSR)
// 	if err != nil {
// 		return "", err
// 	}
// 	return ca.IssueCertificateFromCSR(csr)
// }

// // IssueCertificateFromCSR issues a certificate
// func (ca *CA) IssueCertificateFromCSR(csr *x509.CertificateRequest) (certsPEM string, err error) {
// 	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
// 	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
// 	if err != nil {
// 		return "", err
// 	}

// 	template := &x509.Certificate{
// 		SerialNumber: serialNumber,
// 		Subject:      csr.Subject,
// 		Extensions:   csr.Extensions,
// 		DNSNames:     csr.DNSNames,
// 		NotBefore:    time.Now(),
// 		NotAfter:     ca.certificate.NotAfter,
// 	}
// 	certBuffer, err := x509.CreateCertificate(rand.Reader, template, ca.certificate, csr.PublicKey, ca.privateKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	cert, err := x509.ParseCertificate(certBuffer)
// 	if err != nil {
// 		return "", err
// 	}

// 	return CertToPEM(cert) + CertToPEM(ca.certificate), nil
// }
