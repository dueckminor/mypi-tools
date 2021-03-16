package downloads

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/dueckminor/mypi-tools/go/util"
	"gopkg.in/yaml.v2"
)

func getAlpineMetaData(flavor, arch string, major, minor, patch int) (fm *FileMetadata, err error) {
	fm = &FileMetadata{}
	fm.Name = fmt.Sprintf("%s-%d.%d.%d-%s.tar.gz", flavor, major, minor, patch, arch)
	fm.URL = fmt.Sprintf("https://dl-cdn.alpinelinux.org/alpine/v%d.%d/releases/%s/%s", major, minor, arch, fm.Name)
	fmt.Printf(fm.URL)
	fm.Size, err = downloadGetSize(fm.URL)
	if err != nil {
		return nil, err
	}
	sha256, err := downloadToString(fm.URL + ".sha256")
	if err == nil {
		fm.Checksums = make(map[string]string)
		fm.Checksums["sha256"] = strings.Split(sha256, " ")[0]
	}
	return fm, nil
}

type AlpineDownloader struct {
	Downloader
}

func (d *AlpineDownloader) getAlpineVersions() {
	flavor := "alpine-rpi"
	arch := "aarch64"
	major := 3
	for minor := 13; ; minor++ {
		for patch := 0; ; patch++ {
			fm, err := getAlpineMetaData(flavor, arch, major, minor, patch)
			if err != nil {
				if patch == 0 {
					return
				}
				break
			}
			fm.Tags = map[string]string{
				"alpine-flavor":  flavor,
				"alpine-arch":    arch,
				"alpine-version": fmt.Sprintf("%d.%d.%d", major, minor, patch),
			}
			d.FileMetadatas = append(d.FileMetadatas, fm)
		}
	}
}

func (d *AlpineDownloader) GetDownloaderForVersion(version, arch string) *FileDownloader {
	for _, m := range d.FileMetadatas {
		if m.Tags["alpine-version"] == version && m.Tags["alpine-arch"] == arch {
			return d.GetDownloader(m.Name)
		}
	}
	return nil
}

func NewAlpineDownloader() (d *AlpineDownloader) {
	d = &AlpineDownloader{}

	d.TargetDir = path.Join(os.Getenv("HOME"), ".mypi", "downloads", "alpinelinux")
	os.MkdirAll(d.TargetDir, os.ModePerm)

	versionsYml := path.Join(d.TargetDir, "versions.yml")
	if util.FileExists(versionsYml) {
		data, err := ioutil.ReadFile(versionsYml)
		if err == nil {
			err = yaml.Unmarshal(data, &d.FileMetadatas)
			if err == nil {
				return d
			}
		}
	}

	d.getAlpineVersions()

	if len(d.FileMetadatas) > 0 {
		data, err := yaml.Marshal(d.FileMetadatas)
		if err == nil {
			ioutil.WriteFile(versionsYml, data, os.ModePerm)
		}
	}

	return d
}
