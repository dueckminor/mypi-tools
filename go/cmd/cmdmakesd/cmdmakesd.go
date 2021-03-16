package cmdmakesd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"

	"github.com/dueckminor/mypi-tools/go/cmd"
	"github.com/dueckminor/mypi-tools/go/downloads"
	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/dueckminor/mypi-tools/go/util"
	"github.com/fatih/color"
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
	MypiUUID      string

	BootDevice   string
	RootDevice   string
	WlanSSID     string
	WlanPassword string
	SSH          settingsSSH
}

func extractTarGz(tarFile, destDir string) error {
	f, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if !util.FileIsSafePath(header.Name) {
			continue
		}
		dir, file := path.Split(header.Name)
		toDir := path.Join(destDir, dir)
		toFile := path.Join(toDir, file)

		os.MkdirAll(toDir, os.ModePerm)

		fmt.Println(toFile)

		toFileStream, err := os.Create(toFile)
		defer toFileStream.Close()

		if _, err = io.Copy(toFileStream, tarReader); err != nil {
			return err
		}

		toFileStream.Close()
	}
	return nil
}

func tarCreateFile(tw *tar.Writer, filename string, mode, size int64) error {
	header := new(tar.Header)
	header.Name = filename
	header.Size = size
	header.Mode = mode
	header.ModTime = time.Now()
	// write the header to the tarball archive
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	return nil
}

func tarAddBuffer(tw *tar.Writer, filename string, data []byte, mode int64) error {
	if err := tarCreateFile(tw, filename, mode, int64(len(data))); err != nil {
		return err
	}
	// copy the file data to the tarball
	if _, err := tw.Write(data); err != nil {
		return err
	}
	return nil
}

func tarAddLink(tw *tar.Writer, filename, linkname string, mode int64) error {
	header := new(tar.Header)
	header.Name = filename
	header.Linkname = linkname
	header.Typeflag = tar.TypeSymlink
	header.Mode = int64(mode | int64(os.ModeSymlink))
	header.ModTime = time.Now()
	// write the header to the tarball archive
	return tw.WriteHeader(header)
}

func createAPKOVL(tarfile string, settings *settings) error {
	file, err := os.Create(tarfile)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	staticFiles := path.Join(settings.DirSetup, "static")

	err = filepath.Walk(staticFiles, func(fileName string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relativePath := path.Join(".", fileName[len(staticFiles):])
		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(fileName)
			fmt.Println(relativePath + " -> " + linkTarget)
			if err != nil {
				return err
			}
			tarAddLink(tw, relativePath, linkTarget, int64(info.Mode()))
			return nil
		}
		fmt.Println(relativePath)
		tarCreateFile(tw, relativePath, int64(info.Mode()), info.Size())

		f, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}
		return nil
	})

	templateFiles := path.Join(settings.DirSetup, "templates")
	err = filepath.Walk(templateFiles, func(fileName string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relativePath := path.Join(".", fileName[len(templateFiles):])

		data, err := ioutil.ReadFile(fileName)
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
		tarCreateFile(tw, relativePath, int64(info.Mode()), int64(buf.Len()))

		_, err = buf.WriteTo(tw)
		if err != nil {
			return err
		}

		return nil
	})

	mypiControl := path.Join(settings.DirDist, "mypi-control", "mypi-control-linux-arm64")
	fmt.Println("checking for mypi-control:", mypiControl)
	stat, err := os.Stat(mypiControl)
	if err != nil {
		return err
	}
	tarCreateFile(tw, "mypi-control/bin/mypi-control", 0755, stat.Size())
	f, err := os.Open(mypiControl)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(tw, f)

	return err
}

func (cmd cmdMakeSD) ParseArgs(args []string) (parsedArgs interface{}, err error) {
	settings := &settings{
		Hostname: args[0],
		DirSetup: args[1],
		DirDist:  args[2],
		MypiUUID: uuid.New().String(),
	}
	if len(settings.DirSetup) == 0 {
		settings.DirSetup = *dirSetup
	}
	if len(settings.DirDist) == 0 {
		settings.DirDist = *dirDist
	}
	return settings, nil
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

func createSSHKeys(settings *settings) error {
	// ssh-keygen -t dsa -b 1024 -f ssh_host_dsa_key -N ""
	// ssh-keygen -t rsa -b 3072 -f ssh_host_rsa_key -N ""
	// ssh-keygen -t ecdsa -f ssh_host_ecdsa_key -N ""
	// ssh-keygen -t ed25519 -f ssh_host_ed25519_key -N ""
	return nil
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

	disk, err := fdisk.GetDisk(settings.Disk)
	if err != nil {
		panic(err)
	}

	if !disk.IsRemovable() {
		return nil
	}

	disk.InitializePartitions("MBR", fdisk.PartitionInfo{
		Size:   256 * 1024 * 1024,
		Format: "FAT32",
		Type:   7,
		Name:   "RPI-BOOT",
	})

	fmt.Println("")
	c.Print("                      ")
	fmt.Println("")
	c.Print(" --- Extract data --- ")
	fmt.Println("")
	c.Print("                      ")
	fmt.Println("")
	fmt.Println("")

	disk, err = fdisk.GetDisk(settings.Disk)
	if err != nil {
		panic(err)
	}

	partitions, err := disk.GetPartitions()
	if err != nil {
		panic(err)
	}

	mountPoint, err := partitions[0].GetMountPoint()
	if err != nil {
		panic(err)
	}

	fmt.Println(mountPoint)

	err = extractTarGz(alpineFileDownloader.GetTargetFile(), mountPoint)
	if err != nil {
		panic(err)
	}

	err = createAPKOVL(path.Join(mountPoint, settings.Hostname+".apkovl.tar.gz"), settings)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(path.Join(mountPoint, "mypiuuid.txt"), []byte(settings.MypiUUID+"\n"), os.ModePerm)

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
