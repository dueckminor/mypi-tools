package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func GetFilename(filename string) string {
	if !filepath.IsAbs(filename) {
		if nil == mypiConfig {
			findRoot()
		}
		filename = filepath.Join(mypiRoot, filename)
	}
	return filename
}

func FileToBytes(filename string) ([]byte, error) {
	dat, err := ioutil.ReadFile(GetFilename(filename))
	if err != nil {
		return nil, err
	}
	return dat, err
}

func FileToString(filename string) (string, error) {
	dat, err := FileToBytes(filename)
	if err != nil {
		return "", err
	}
	return string(dat), err
}

func ReadYAML(filename string, result interface{}) error {
	dat, err := FileToBytes(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(dat, result)
}

func ReadJSON(filename string, result interface{}) error {
	dat, err := FileToBytes(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(dat, result)
}
