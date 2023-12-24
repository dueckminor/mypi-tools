package webhandler

import (
	"net/http"

	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/gin-gonic/gin"
)

// GetDisks return list of hosts
func GetDisks(c *gin.Context) {
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
	c.JSON(http.StatusOK, diskInfos)
}
