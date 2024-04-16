package cmdmakesd

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/dueckminor/mypi-tools/go/cmd"
	"github.com/dueckminor/mypi-tools/go/downloads"
	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/dueckminor/mypi-tools/go/util/panic"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"

	. "github.com/dueckminor/mypi-tools/go/setup"
)

var (
	dirSetup = flag.String("dir-setup", "", "the directory containing the setup files")
	dirDist  = flag.String("dir-dist", "", "the directory containing the dist files")
)

type cmdMakeSD struct{}

type settingsSSH struct {
	AuthorizedKeys string
}

type settings struct {
	Disk          string
	Hostname      string
	AlpineVersion string
	AlpineArch    string
	DirSetup      string
	DirDist       string
	DirTarget     string
	ZipTarget     string
	MypiUUID      string

	BootDevice   string
	RootDevice   string
	WlanSSID     string
	WlanPassword string
	SSH          settingsSSH
}

func writeSimpleFile(w DirWriter, filename string, content []byte) error {

	f, err := w.CreateFile(FileInfo{Type: FileTypeFile, Name: filename})
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return nil
}

type staticFileInfos struct {
	Executables []string
	Softlinks   map[string]string
}

func (s *staticFileInfos) GetLinkInfo(fileName string) (linkTarget string, ok bool) {
	fileName = strings.TrimPrefix(fileName, "./")
	if linkTarget, ok = s.Softlinks[fileName]; ok {
		return linkTarget, ok
	}
	return "", false
}

func (s *staticFileInfos) GetFileMode(fileName string) (fileMode int64, ok bool) {
	fileName = strings.TrimPrefix(fileName, "./")
	if fileName == "fileinfos.yml" {
		return 0, false
	}

	parts := strings.Split(fileName, "/")

	for _, pattern := range s.Executables {
		patternParts := strings.Split(pattern, "/")
		if len(patternParts) != len(parts) {
			continue
		}
		match := true
		for i, part := range parts {
			if patternParts[i] == "*" || patternParts[i] == part {
				continue
			}
			match = false
			break
		}
		if match {
			return int64(0755), true
		}

	}

	return int64(0644), true
}

func createAPKOVL(w DirWriter, filename string, settings *settings) error {
	file, err := w.CreateFile(FileInfo{Type: FileTypeFile, Name: filename})
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw, err := NewTarWriter(gw)
	if err != nil {
		return err
	}
	defer tw.Close()

	staticFiles := path.Join(settings.DirSetup, "static")

	staticFileInfosYml, err := os.ReadFile(path.Join(staticFiles, "fileinfos.yml"))
	if err != nil {
		return err
	}
	var staticFileInfos staticFileInfos
	err = yaml.Unmarshal(staticFileInfosYml, &staticFileInfos)
	if err != nil {
		return err
	}

	err = filepath.Walk(staticFiles, func(fileName string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relativePath := path.Join(".", fileName[len(staticFiles):])
		if linkTarget, ok := staticFileInfos.GetLinkInfo(relativePath); ok {
			return tw.AddLink(relativePath, linkTarget, int64(0644))
		} else if fileMode, ok := staticFileInfos.GetFileMode(relativePath); ok {
			w, err := tw.CreateFile(relativePath, fileMode, info.Size())
			if err != nil {
				return err
			}
			f, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(w, f)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	templateFiles := path.Join(settings.DirSetup, "templates")
	err = filepath.Walk(templateFiles, func(fileName string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relativePath := path.Join(".", fileName[len(templateFiles):])

		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		t, err := template.New(relativePath).Parse(string(data))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, settings)
		if err != nil {
			return err
		}

		fmt.Println(relativePath)
		w, err := tw.CreateFile(relativePath, int64(info.Mode()), int64(buf.Len()))
		if err != nil {
			return err
		}

		_, err = buf.WriteTo(w)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	mypiControl := path.Join(settings.DirDist, "mypi-control", "mypi-control-linux-arm64")
	fmt.Println("checking for mypi-control:", mypiControl)
	stat, err := os.Stat(mypiControl)
	if err != nil {
		return err
	}
	w2, err := tw.CreateFile("mypi-control/bin/mypi-control", 0755, stat.Size())
	if err != nil {
		return err
	}
	f, err := os.Open(mypiControl)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w2, f)

	return err
}

func (cmd cmdMakeSD) ParseArgs(args []string) (parsedArgs interface{}, err error) {
	settings := &settings{}

	f := flag.NewFlagSet("makesd", flag.ExitOnError)
	f.StringVar(&settings.AlpineVersion, "alpine-version", "latest", "")
	f.StringVar(&settings.AlpineArch, "alpine-arch", "aarch64", "")
	f.StringVar(&settings.Disk, "disk", "", "")
	f.StringVar(&settings.Hostname, "hostname", "", "")
	f.StringVar(&settings.DirSetup, "dir-setup", "", "")
	f.StringVar(&settings.DirDist, "dir-dist", "", "")
	f.StringVar(&settings.DirTarget, "dir-target", "", "")
	f.StringVar(&settings.ZipTarget, "zip-target", "", "")

	err = f.Parse(args)
	if err != nil {
		return nil, err
	}

	return settings, err
}

func (cmd cmdMakeSD) UnmarshalArgs(marshaledArgs []byte) (parsedArgs interface{}, err error) {
	settings := &settings{}
	err = json.Unmarshal(marshaledArgs, &settings)
	if len(settings.DirSetup) == 0 {
		settings.DirSetup = *dirSetup
	}
	if len(settings.DirDist) == 0 {
		settings.DirDist = *dirDist
	}
	return settings, err
}

func (cmd cmdMakeSD) Execute(parsedArgs interface{}) error {
	settings, ok := parsedArgs.(*settings)
	if !ok {
		return os.ErrInvalid
	}

	c := color.New(color.BgBlue).Add(color.FgHiYellow)

	alpineDownloader := downloads.NewAlpineDownloader()
	alpineFileDownloader := alpineDownloader.GetDownloaderForVersion(settings.AlpineVersion, settings.AlpineArch)

	alpineFileDownloader.StartDownload()
	if !alpineFileDownloader.DownloadCompleted() {
		fmt.Println("")
		c.Print("                                  ")
		fmt.Println("")
		c.Print(" --- Downloading Alpine Linux --- ")
		fmt.Println("")
		c.Print("                                  ")
		fmt.Println("")
		fmt.Println("")

		for !alpineFileDownloader.DownloadCompleted() {
			fmt.Print(".")
			time.Sleep(5 * time.Second)
		}
		fmt.Println("")
		fmt.Println("completed")
	}

	fmt.Println("FileName:", alpineFileDownloader.MetaData.Name)

	fmt.Println("")
	c.Print("                           ")
	fmt.Println("")
	c.Print(" --- Initializing disk --- ")
	fmt.Println("")
	c.Print("                           ")
	fmt.Println("")
	fmt.Println("")

	fmt.Println("dir-dist", settings.DirDist)
	fmt.Println("dir-setup", settings.DirSetup)

	var w DirWriter
	var err error

	dirTarget := settings.DirTarget
	if len(settings.ZipTarget) > 0 {
		filewriter, err := os.Create(settings.ZipTarget)
		panic.OnError(err)
		w, err = NewZipWriter(filewriter)
		panic.OnError(err)
	} else {
		if len(dirTarget) == 0 {
			disk, err := fdisk.GetDisk(settings.Disk)
			panic.OnError(err)

			if !disk.IsRemovable() {
				return nil
			}

			err = disk.InitializePartitions("MBR", fdisk.PartitionInfo{
				Size:   256 * 1024 * 1024,
				Format: "FAT32",
				Type:   7,
				Name:   "RPI-BOOT",
			})
			panic.OnError(err)

			disk, err = fdisk.GetDisk(settings.Disk)
			panic.OnError(err)

			partitions, err := disk.GetPartitions()
			panic.OnError(err)

			dirTarget, err = partitions[0].GetMountPoint()
			panic.OnError(err)
		}
		w, err = NewFileWriter(dirTarget)
		panic.OnError(err)
	}
	defer w.Close()

	err = TarGzFileExtract(alpineFileDownloader.GetTargetFile(), w)
	panic.OnError(err)

	// there has to be exactly ONE `*.apkovl.tar.gz` in the root-directory
	err = createAPKOVL(w, settings.Hostname+".apkovl.tar.gz", settings)
	panic.OnError(err)

	err = writeSimpleFile(w, "mypiuuid.txt", []byte(settings.MypiUUID+"\n"))
	panic.OnError(err)

	fmt.Println("")
	c.Print("                   ")
	fmt.Println("")
	c.Print(" --- Succeeded --- ")
	fmt.Println("")
	c.Print("                   ")
	fmt.Println("")
	fmt.Println("")

	return nil
}

func init() {
	cmd.Register("makesd", cmdMakeSD{})
}
