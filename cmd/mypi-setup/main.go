package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-gonic/gin"

	"github.com/dueckminor/mypi-tools/go/gotty/localcommand"
	"github.com/dueckminor/mypi-tools/go/gotty/server/ginhandler"

	"github.com/dueckminor/mypi-tools/go/cmd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdmakesd"
)

func setupDisk() error {
	disks, err := fdisk.GetDisks()
	if err != nil {
		return err
	}
	for _, disk := range disks {
		if disk.IsRemovable() {
			fmt.Println("Removeable-Disk: " + disk.GetDeviceName())
			disk.InitializePartitions("MBR", fdisk.PartitionInfo{
				Size:   256 * 1024 * 1024,
				Format: "FAT32",
				Type:   7,
				Name:   "RPI-BOOT",
			})
		}
	}
	return nil
}

func main() {
	if len(os.Args) > 1 {
		if cmd.IsAvailable(os.Args[1]) {
			cmd.Execute(os.Args[1], os.Args[2:]...)
			return
		}
	}

	flag.Parse()

	r := gin.Default()

	wh := &webhandler.WebHandler{}
	err := wh.SetupEndpoints(r)

	r.GET("/api/disks", func(c *gin.Context) {
		type DiskInfo struct {
			Name string `json:"name,omitempty"`
			Size int64  `json:"size"`
		}
		diskInfos := []DiskInfo{}
		disks, err := fdisk.GetDisks()
		if err == nil {
			for _, disk := range disks {
				if disk.IsRemovable() {
					diskInfos = append(diskInfos, DiskInfo{
						Name: disk.GetDeviceName(),
						Size: disk.GetSize(),
					})
				}
			}
		}

		data, err := json.Marshal(diskInfos)
		c.Data(200, "application/json", data)
	})
	r.POST("/api/actions/makesd", func(c *gin.Context) {
		type MakeSD struct {
			DiskName      string `json:"diskname,omitempty"`
			AlpineVersion string `json:"alpineversion,omitempty"`
			AlpineArch    string `json:"alpinearch,omitempty"`
			HostName      string `json:"hostname,omitempty"`
		}

	})

	r.GET("/ws", func(c *gin.Context) {
		factory, err := localcommand.NewFactory(os.Args[0], []string{"makesd"}, &localcommand.Options{})
		if err == nil {
			ginhandler.Handler(c, factory)
		}
	})

	if err != nil {
		panic(err)
	}

	restapi.Run(r)
}
