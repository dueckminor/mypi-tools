package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return true
	}
	return false
}

var (
	mypiRoot   string
	mypiConfig Config
)

func findRoot() string {
	if 0 != len(mypiRoot) {
		return mypiRoot
	}

	var err error

	mypiYML := os.Getenv("MYPI_CONFIG")
	mypiRoot = os.Getenv("MYPI_ROOT")

	for mypiYML != mypiRoot+"/config/mypi.yml" {
		if len(mypiYML) > 0 && fileExists(mypiYML) {
			mypiConfig, err = ReadConfigFile(mypiYML)
			if err != nil {
				panic(err)
			}
			root := mypiConfig.GetString("config", "root")
			if len(root) > 0 {
				mypiRoot = root
				mypiYML = mypiRoot + "/config/mypi.yml"
				mypiConfig = nil
			}
		} else if len(mypiRoot) > 0 {
			mypiYML = mypiRoot + "/config/mypi.yml"
		} else {
			dir, _ := os.Executable()
			dir = filepath.Dir(dir)
			mypiYML = filepath.Join(dir, "mypi.yml")
			if !fileExists(mypiYML) {
				dir, _ := os.Getwd()
				mypiYML = filepath.Join(dir, "mypi.yml")
			}
			if !fileExists(mypiYML) {
				panic("cant find mypi.yml")
			}
		}
	}
	if nil == mypiConfig {
		mypiConfig, err = ReadConfigFile(mypiYML)
	}
	return mypiRoot
}

type Config interface {
	GetString(path ...interface{}) string
	GetBool(path ...interface{}) bool
}

type configImpl struct {
	cfg interface{}
}

func (c configImpl) FollowPath(path ...interface{}) (interface{}, error) {
	cfg := c.cfg
	for _, pathElement := range path {
		switch pathElementV := pathElement.(type) {
		case int:
			a, ok := cfg.([]interface{})
			if !ok {
				return nil, nil
			}
			if (pathElementV < 0) || (pathElementV >= len(a)) {
				return nil, nil
			}
			cfg = a[pathElementV]
		case string:
			m, ok := cfg.(map[interface{}]interface{})
			if !ok {
				return nil, nil
			}
			cfg, ok = m[pathElementV]
			if !ok {
				return nil, nil
			}
		default:
			return nil, nil
		}
	}
	return cfg, nil
}

func (c configImpl) GetString(path ...interface{}) string {
	cfg, err := c.FollowPath(path...)
	if err != nil {
		return ""
	}
	switch v := cfg.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func (c configImpl) GetBool(path ...interface{}) bool {
	cfg, err := c.FollowPath(path...)
	if err != nil {
		return false
	}
	switch v := cfg.(type) {
	case bool:
		return v
	case int:
		return (v != 0)
	default:
		return false
	}
	return false
}

func ReadConfigFile(filename string) (Config, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	r := configImpl{}
	err = yaml.Unmarshal(dat, &r.cfg)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func GetConfig() (c Config) {
	if nil == mypiConfig {
		findRoot()
	}
	return mypiConfig
}

func GetRoot() string {
	if len(mypiRoot) == 0 {
		findRoot()
	}
	return mypiRoot
}
