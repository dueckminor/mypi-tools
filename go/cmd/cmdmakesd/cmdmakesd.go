package cmdmakesd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/dueckminor/mypi-tools/go/cmd"
	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/fatih/color"
)

type cmdMakeSD struct{}

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

func (cmd cmdMakeSD) Execute(args []string) error {

	c := color.New(color.BgBlue).Add(color.FgHiYellow)

	fmt.Println("")
	c.Print("                           ")
	fmt.Println("")
	c.Print(" --- Initializing disk --- ")
	fmt.Println("")
	c.Print("                           ")
	fmt.Println("")
	fmt.Println("")

	disk, err := fdisk.GetDisk("disk5")
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

	disk, err = fdisk.GetDisk("disk5")
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

	extractTarGz(path.Join(os.Getenv("HOME"), "Downloads", "alpine-rpi-3.11.3-aarch64.tar.gz"), mountPoint)

	return nil
}

func init() {
	cmd.Register("makesd", cmdMakeSD{})
}
