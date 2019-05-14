package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ReadRSAPrivateKey(filename string) (pk *rsa.PrivateKey, err error) {
	blob, err := FileToBytes(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(blob)

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pk, ok := parseResult.(*rsa.PrivateKey)
		if ok {
			return pk, nil
		}
	}
	return nil, errors.New("No RSA Key found")
}

func BlobToRSAPublicKey(blob []byte) (pk *rsa.PublicKey, err error) {
	block, _ := pem.Decode(blob)

	switch block.Type {
	case "PUBLIC KEY":
		parseResult, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pk, ok := parseResult.(*rsa.PublicKey)
		if ok {
			return pk, nil
		}
	}
	return nil, errors.New("No RSA Key found")
}

func StringToRSAPublicKey(str string) (pk *rsa.PublicKey, err error) {
	return BlobToRSAPublicKey([]byte(str))
}

func ReadRSAPublicKey(filename string) (pk *rsa.PublicKey, err error) {
	blob, err := FileToBytes(filename)
	if err != nil {
		return nil, err
	}
	return BlobToRSAPublicKey(blob)
}
