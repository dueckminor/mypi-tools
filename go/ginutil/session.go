package ginutil

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func createKeyPair() (auth []byte, enc []byte, err error) {
	auth = make([]byte, 64)
	_, err = rand.Read(auth)
	if err != nil {
		return nil, nil, err
	}
	enc = make([]byte, 32)
	_, err = rand.Read(enc)
	if err != nil {
		return nil, nil, err
	}
	return auth, enc, nil
}

func storeKey(key []byte, cfg config.Config, name string) (err error) {
	return cfg.SetString(name, base64.StdEncoding.EncodeToString(key))
}
func storeKeyPair(auth []byte, enc []byte, cfg config.Config) (err error) {
	err = storeKey(auth, cfg, "auth")
	if err != nil {
		return err
	}
	return storeKey(enc, cfg, "enc")
}

func readKey(cfg config.Config, name string) []byte {
	val := cfg.GetString(name)
	if len(val) == 0 {
		return nil
	}
	key, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil
	}
	return key
}
func readKeyPair(cfg config.Config) (auth []byte, enc []byte) {
	return readKey(cfg, "auth"), readKey(cfg, "enc")
}

func ConfigureSessionCookies(r *gin.Engine, cfg config.Config, session string) (err error) {
	cfgSession, err := cfg.MakeMap("session")
	if err != nil {
		return err
	}
	cfgCurrent, err := cfgSession.MakeMap("current")
	if err != nil {
		return err
	}

	auth, enc := readKeyPair(cfgCurrent)
	if auth == nil {
		auth, enc, err = createKeyPair()
		if err != nil {
			return err
		}
		err = storeKeyPair(auth, enc, cfgCurrent)
		if err != nil {
			return err
		}
		err = cfg.Write()
		if err != nil {
			return err
		}
	}

	store := cookie.NewStore(auth, enc)
	r.Use(sessions.Sessions(session, store))

	return nil
}
