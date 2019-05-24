package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

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

	for i, arg := range os.Args[1:] {

		dir := path.Join(tmpDir, fmt.Sprintf("%d", i))
		path := fmt.Sprintf("/cams/%d/", i)

		os.MkdirAll(dir, os.ModePerm)

		go runFFMPEG(dir, arg)

		r.Use(static.ServeRoot(path, dir))
		r.Use(static.ServeRoot(path, "./web"))
	}

	panic(r.Run(":8080"))
}
