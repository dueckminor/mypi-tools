package webhandler

import (
	"github.com/dueckminor/mypi-tools/go/downloads"
	"github.com/gin-gonic/gin"
)

type WebhandlerDownloads struct {
	downloader *downloads.AlpineDownloader
}

func NewWebhandlerDownloads() (wh *WebhandlerDownloads) {
	wh = &WebhandlerDownloads{}
	wh.downloader = downloads.NewAlpineDownloader()
	return wh
}

func (wh *WebhandlerDownloads) getDownloads(c *gin.Context) {
	c.JSON(200, wh.downloader.GetMetadatas())
}

func (wh *WebhandlerDownloads) SetupEndpoints(r *gin.RouterGroup) {
	r.GET("/api/downloads/alpine", wh.getDownloads)
}
