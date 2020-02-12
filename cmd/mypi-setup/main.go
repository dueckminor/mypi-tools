package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-gonic/gin"
	//	"github.com/gorilla/websocket"
)

// var wsupgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func wshandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := wsupgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		//fmt.Println("Failed to set websocket upgrade: %+v", err)
// 		return
// 	}

// 	for {
// 		t, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		conn.WriteMessage(t, msg)
// 	}
// }

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
	r.POST("/api/actions/create_sd", func(c *gin.Context) {
		type CreateSD struct {
			DiskName      string `json:"diskname,omitempty"`
			AlpineVersion string `json:"alpineversion,omitempty"`
			AlpineArch    string `json:"alpinearch,omitempty"`
			HostName      string `json:"hostname,omitempty"`
		}

	})

	// r.GET("/ws", func(c *gin.Context) {
	// 	wshandler(c.Writer, c.Request)
	// })

	if err != nil {
		panic(err)
	}

	restapi.Run(r)
}
