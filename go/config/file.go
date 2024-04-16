package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func GetFilename(filenameParts ...string) string {
	if len(filenameParts) == 0 || !filepath.IsAbs(filenameParts[0]) {
		if nil == mypiConfig {
			findRoot()
		}
		parts := []string{mypiRoot}
		parts = append(parts, filenameParts...)
		return filepath.Join(parts...)
	}
	return filepath.Join(filenameParts...)
}

func FileToBytes(filenameParts ...string) ([]byte, error) {
	dat, err := os.ReadFile(GetFilename(filenameParts...))
	if err != nil {
		return nil, err
	}
	return dat, err
}

func FileToString(filenameParts ...string) (string, error) {
	dat, err := FileToBytes(filenameParts...)
	if err != nil {
		return "", err
	}
	return string(dat), err
}

func ReadYAML(result interface{}, filenameParts ...string) error {
	dat, err := FileToBytes(filenameParts...)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(dat, result)
}

func ReadJSON(result interface{}, filenameParts ...string) error {
	dat, err := FileToBytes(filenameParts...)
	if err != nil {
		return err
	}
	return json.Unmarshal(dat, result)
}
