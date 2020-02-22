package cmdmakesd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/dueckminor/mypi-tools/go/cmd"
	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/fatih/color"
)

type cmdMakeSD struct{}

type settings struct {
	Hostname string
	DirSetup string
	DirDist  string

	BootDevice   string
	RootDevice   string
	WlanSSID     string
	WlanPassword string
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

func createAPKOVL(tarfile string, settings settings) error {
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

		tarCreateFile(tw, relativePath, int64(info.Mode()), int64(buf.Len()))

		_, err = buf.WriteTo(tw)
		if err != nil {
			return err
		}

		return nil
	})

	mypiSetup := path.Join(settings.DirDist, "mypi-setup-linux-arm64")
	stat, err := os.Stat(mypiSetup)
	if err != nil {
		return err
	}
	tarCreateFile(tw, "mypi-setup/bin/mypi-setup", 0755, stat.Size())
	f, err := os.Open(mypiSetup)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(tw, f)

	return err
}

func (cmd cmdMakeSD) Execute(args []string) error {

	c := color.New(color.BgBlue).Add(color.FgHiYellow)

	settings := settings{
		Hostname: args[0],
		DirSetup: args[1],
		DirDist:  args[2],
	}

	fmt.Println("")
	c.Print("                           ")
	fmt.Println("")
	c.Print(" --- Initializing disk --- ")
	fmt.Println("")
	c.Print("                           ")
	fmt.Println("")
	fmt.Println("")

	disk, err := fdisk.GetDisk("/dev/disk2")
	if err != nil {
		panic(err)
	}

	if !disk.IsRemovable() {
		return nil
	}

	// disk.InitializePartitions("MBR", fdisk.PartitionInfo{
	// 	Size:   256 * 1024 * 1024,
	// 	Format: "FAT32",
	// 	Type:   7,
	// 	Name:   "RPI-BOOT",
	// })

	fmt.Println("")
	c.Print("                      ")
	fmt.Println("")
	c.Print(" --- Extract data --- ")
	fmt.Println("")
	c.Print("                      ")
	fmt.Println("")
	fmt.Println("")

	disk, err = fdisk.GetDisk("/dev/disk2")
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

	//extractTarGz(path.Join(os.Getenv("HOME"), "Downloads", "alpine-rpi-3.11.3-aarch64.tar.gz"), mountPoint)

	createAPKOVL(path.Join(mountPoint, settings.Hostname+".apkovl.tar.gz"), settings)

	return nil
}

func init() {
	cmd.Register("makesd", cmdMakeSD{})
}
