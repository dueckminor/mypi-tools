package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	mypiRoot = flag.String("mypi-root", "", "The root of the mypi filesystem")
)

func init() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		config.InitApp(*mypiRoot)
	}
}

func runFFMPEG(dir, url string) {
	for {
		args := []string{"-i", url}

		args = append(args, "-an")
		args = append(args, "-c:v")
		args = append(args, "copy")

		args = append(args, "-f")
		args = append(args, "hls")
		args = append(args, "-hls_flags")
		args = append(args, "delete_segments")
		args = append(args, "-hls_list_size")
		args = append(args, "4")

		args = append(args, path.Join(dir, "stream.m3u8"))

		// args = append(args, "-an")
		// args = append(args, "-c:v")
		// args = append(args, "copy")
		// args = append(args, "-f")
		// args = append(args, "dash")
		// args = append(args, "-window_size")
		// args = append(args, "4")
		// args = append(args, path.Join(dir, "stream.mpd"))

		cmd := exec.Command("ffmpeg", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		cmd.Wait()
	}
}

func main() {
	r := gin.Default()

	tmpDir := os.TempDir()
	if _, err := os.Stat("/dev/shm"); !os.IsNotExist(err) {
		tmpDir = "/dev/shm"
	}

	tmpDir = path.Join(tmpDir, "videostream")
	os.MkdirAll(tmpDir, os.ModePerm)

	fmt.Println("Using dir:", tmpDir)

	iCam := 0

	for _, url := range flag.Args() {

		dir := path.Join(tmpDir, fmt.Sprintf("%d", iCam))
		path := fmt.Sprintf("/cams/%d/", iCam)

		os.MkdirAll(dir, os.ModePerm)

		go runFFMPEG(dir, url)

		r.GET("/cams/", func(c *gin.Context) {
			c.Header("Expires", "0")
		}, static.ServeRoot(path, dir))

		r.Use(static.ServeRoot(path, "./web"))
		iCam++
	}

	for _, e := range config.GetConfig().GetArray("config", "webcams") {
		url := e.GetString("url")

		dir := path.Join(tmpDir, fmt.Sprintf("%d", iCam))
		path := fmt.Sprintf("/cams/%d/", iCam)

		os.MkdirAll(dir, os.ModePerm)

		go runFFMPEG(dir, url)

		r.GET("/cams/", func(c *gin.Context) {
			c.Header("Expires", "0")
		}, static.ServeRoot(path, dir))

		r.Use(static.ServeRoot(path, "./web"))
		iCam++
	}

	restapi.Run(r)
}
