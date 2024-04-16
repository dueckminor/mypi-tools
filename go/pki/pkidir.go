package pki

import (
	"os"
	"path"
	"strings"
)

// PKIDir provides access to certificate/private_keys stored in a directory
type PKIDir interface {
	GetGenerator() *PkiGenerator
	GetIdentities() (identities []Identity, err error)
}

type pkiDir struct {
	dir string
}

func NewPKIDir(dir string) (p PKIDir, err error) {
	return &pkiDir{dir: dir}, nil
}

func (p *pkiDir) GetGenerator() *PkiGenerator {
	return &PkiGenerator{
		PkiDir: p.dir,
	}
}

func (p *pkiDir) GetIdentities() (identities []Identity, err error) {
	files, err := os.ReadDir(p.dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, "_priv.pem") {
			identity, err := LoadIdentity(path.Join(p.dir, name[:len(name)-9]))
			if err == nil {
				identities = append(identities, identity)
			}
		}
	}
	return identities, nil
}
