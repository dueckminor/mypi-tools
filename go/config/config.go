package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dueckminor/mypi-tools/go/util"
	yaml "gopkg.in/yaml.v3"
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

	if len(mypiRoot) == 0 {
		mypiRoot = "/opt/mypi"
	}

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
	Get(path ...interface{}) (Config, error)
	GetString(path ...interface{}) string
	GetBool(path ...interface{}) bool
	GetArray(path ...interface{}) []Config
	GetMap(path ...interface{}) (ConfigMap, error)
	MakeMap(path ...interface{}) (ConfigMap, error)
	AddArrayElement(obj interface{}, path ...interface{}) error
	SetString(name, value string) error
	Write() error
}

type ConfigMap interface {
	Config
}

type configImpl struct {
	cfg interface{}
}

type configMapImpl struct {
	configImpl
}

type configImplOnFile struct {
	configImpl
	filename string
}

func toMap(cfg interface{}) map[string]interface{} {
	if m, ok := cfg.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func toArray(cfg interface{}) []interface{} {
	if a, ok := cfg.([]interface{}); ok {
		return a
	}
	return nil
}

func followPath(cfg interface{}, path ...interface{}) (interface{}, error) {
	for _, pathElement := range path {
		switch pathElementV := pathElement.(type) {
		case int:
			a := toArray(cfg)
			if a == nil {
				return nil, nil
			}
			if (pathElementV < 0) || (pathElementV >= len(a)) {
				return nil, nil
			}
			cfg = a[pathElementV]
		case string:
			m := toMap(cfg)
			if m == nil {
				return nil, nil
			}
			var ok bool
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

func (c *configImpl) FollowPath(path ...interface{}) (interface{}, error) {
	return followPath(c.cfg, path...)
}

func (c *configImpl) Get(path ...interface{}) (Config, error) {
	if len(path) == 0 {
		return c, nil
	}
	var err error
	result := &configImpl{}
	result.cfg, err = followPath(c.cfg, path...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *configImpl) GetMap(path ...interface{}) (ConfigMap, error) {
	return nil, nil
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
	case map[string]interface{}:
		v[name] = value
		return nil
	}
	return errors.New("wrong type")
}

func (c *configImpl) MakeMap(path ...interface{}) (result ConfigMap, err error) {
	if len(path) == 0 {
		if c.cfg != nil {
			if nil == toMap(c.cfg) {
				return nil, errors.New("wrong type")
			}
		} else {
			c.cfg = make(map[string]interface{})
		}
		return &configMapImpl{configImpl: *c}, nil
	}
	switch v := path[0].(type) {
	case string:
		if c.cfg == nil {
			c.cfg = make(map[string]interface{})
		}
		m := toMap(c.cfg)
		if m != nil {
			if len(path) == 1 {
				if e, ok := m[v]; ok {
					m = toMap(e)
					if m != nil {
						return &configMapImpl{configImpl: configImpl{cfg: m}}, nil
					} else {
						return nil, errors.New("wrong type")
					}
				}
				m[v] = make(map[string]interface{})
				return &configMapImpl{configImpl: configImpl{cfg: m[v]}}, nil
			} else {
				panic("not yet implemented")
			}
		}
	}
	panic("not yet implemented")
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
	r := &configImplOnFile{filename: GetFilename(filename)}
	err = yaml.Unmarshal(dat, &r.cfg)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func ReadConfigFile(filename string) (Config, error) {
	return readConfigFile(GetFilename(filename))
}

func GetOrCreateConfigFile(filename string) (Config, error) {
	filename = GetFilename(filename)
	if util.FileExists(filename) {
		return readConfigFile(filename)
	}
	return New(filename, map[string]interface{}{}), nil
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

// InitApp creates an in-memory-config containing only a single value:
// `config.root`
func InitApp(root string) (err error) {
	mypiRoot = root

	mypiConfigFile := filepath.Join(mypiRoot, "config/mypi.yml")
	if fileExists(mypiConfigFile) {
		mypiConfig, err = readConfigFile(mypiConfigFile)
		if err == nil {
			return nil
		}
	}

	mypiConfig = &configImpl{}

	cfg, err := mypiConfig.MakeMap("config")
	if err != nil {
		return err
	}
	cfg.SetString("root", root)

	fmt.Println(mypiConfig.GetString("config", "root"))

	return nil
}
