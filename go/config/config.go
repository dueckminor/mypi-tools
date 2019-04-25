package config

import (
	"errors"
	"fmt"
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
			mypiConfig, err = readConfigFile(mypiYML)
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
		mypiConfig, err = readConfigFile(mypiYML)
	}
	return mypiRoot
}

type Config interface {
	GetString(path ...interface{}) string
	GetBool(path ...interface{}) bool
	GetArray(path ...interface{}) []Config
	AddArrayElement(obj interface{}, path ...interface{}) error
	SetString(name, value string) error
	Write() error
}

type configImpl struct {
	cfg interface{}
}

type configImplOnFile struct {
	configImpl
	filename string
}

func (c *configImpl) FollowPath(path ...interface{}) (interface{}, error) {
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

func (c *configImpl) GetString(path ...interface{}) string {
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

func (c *configImpl) SetString(name, value string) error {
	switch v := c.cfg.(type) {
	case map[interface{}]interface{}:
		v[name] = value
		return nil
	}
	return errors.New("wrong type")
}

func (c *configImpl) GetBool(path ...interface{}) bool {
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

func (c *configImpl) GetArray(path ...interface{}) []Config {
	cfg, err := c.FollowPath(path...)
	if err != nil {
		return nil
	}
	a, ok := cfg.([]interface{})
	if !ok {
		return nil
	}
	result := make([]Config, len(a))
	for i, e := range a {
		result[i] = &configImpl{cfg: e}
	}
	return result
}

func (c *configImpl) AddArrayElement(obj interface{}, path ...interface{}) error {
	cfg, err := c.FollowPath(path...)
	if err != nil {
		return err
	}
	a, ok := cfg.([]interface{})
	if !ok {
		fmt.Println(cfg)
		if cfg != nil {
			return nil
		}
		a = make([]interface{}, 1)
	}
	a = append(a, obj)

	c.cfg = a
	return nil
}

func (c *configImpl) Write() error {
	return nil
}

func (c *configImplOnFile) Write() error {

	data, err := yaml.Marshal(c.cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.filename, data, 0644)
}

func readConfigFile(filename string) (Config, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	r := &configImplOnFile{filename: filename}
	err = yaml.Unmarshal(dat, &r.cfg)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func ReadConfigFile(filename string) (Config, error) {
	if nil == mypiConfig {
		findRoot()
	}

	if !filepath.IsAbs(filename) {
		filename = filepath.Join(mypiRoot, filename)
	}
	return readConfigFile(filename)
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

func New(filename string, cfg interface{}) (c Config) {
	if nil == mypiConfig {
		findRoot()
	}
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(mypiRoot, filename)
	}
	return &configImplOnFile{
		configImpl: configImpl{cfg: cfg},
		filename:   filename,
	}
}
