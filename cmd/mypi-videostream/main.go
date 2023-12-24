package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	// provide only the http rest api
	_ "github.com/dueckminor/mypi-tools/go/restapi/http"
)

var (
	mypiRoot = flag.String("mypi-root", "", "The root of the mypi filesystem")
)

func init() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		err := config.InitApp(*mypiRoot)
		if err != nil {
			panic(err)
		}
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
		err := cmd.Start()
		if err != nil {
			fmt.Println("ffmpeg failed to start:", err)
			time.Sleep(time.Second * 30)
			continue
		}
		cmd.Wait() // nolint: errcheck
	}
}

func main() {
	r := gin.Default()

	tmpDir := os.TempDir()
	if _, err := os.Stat("/dev/shm"); !os.IsNotExist(err) {
		tmpDir = "/dev/shm"
	}

	tmpDir = path.Join(tmpDir, "videostream")
	err := os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("Using dir:", tmpDir)

	iCam := 0

	addCam := func(url string) {
		dir := path.Join(tmpDir, fmt.Sprintf("%d", iCam))
		path := fmt.Sprintf("/cams/%d/", iCam)
		filter := fmt.Sprintf("%s/:file", path)

		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		go runFFMPEG(dir, url)

		r.GET(filter, func(c *gin.Context) {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}, static.ServeRoot(path, dir))

		iCam++
	}

	for _, url := range flag.Args() {
		addCam(url)
	}

	for _, e := range config.GetConfig().GetArray("config", "webcams") {
		addCam(e.GetString("url"))
	}

	restapi.Run(r)
}
