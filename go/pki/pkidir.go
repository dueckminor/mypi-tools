package pki

// PKIDir provides access to certificate/private_keys stored in a directory
type PKIDir interface {
}

type pkiDir struct {
	dir string
}

func NewPKIDir(dir string) (p PKIDir, err error) {
	return &pkiDir{dir: dir}, nil
}
