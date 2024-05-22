package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/dueckminor/mypi-tools/go/config"
)

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.

	rootCaCert := config.GetFilename("config/pki/root_ca_cert.pem")

	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile(rootCaCert)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}
	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(config.GetFilename("config/pki/tls_cert.pem"), config.GetFilename("config/pki/tls_priv.pem"))
	if err != nil {
		panic(err)
	}
	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}
