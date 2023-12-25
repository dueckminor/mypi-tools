package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

type FFMPEG struct {
	dir  string
	url  string
	file string
	cmd  *exec.Cmd
}

func (ff *FFMPEG) makeArgs() []string {
	args := []string{"-i", ff.url}

	args = append(args, "-an")
	args = append(args, "-c:v")
	args = append(args, "copy")

	args = append(args, "-f")
	args = append(args, "hls")
	args = append(args, "-hls_flags")
	args = append(args, "delete_segments")
	args = append(args, "-hls_list_size")
	args = append(args, "4")

	ff.file = path.Join(ff.dir, "stream.m3u8")
	args = append(args, ff.file)

	return args
}

func (ff *FFMPEG) makeArgsDash() []string {
	args := []string{"-i", ff.url}

	args = append(args, "-an")
	args = append(args, "-c:v")
	args = append(args, "copy")

	args = append(args, "-an")
	args = append(args, "-c:v")
	args = append(args, "copy")
	args = append(args, "-f")
	args = append(args, "dash")
	args = append(args, "-window_size")
	args = append(args, "4")

	ff.file = path.Join(ff.dir, "stream.mpd")
	args = append(args, ff.file)

	return args
}

func (ff *FFMPEG) exec(ctx context.Context) (done chan bool, err error) {
	if ff.cmd != nil {
		return nil, nil
	}
	args := ff.makeArgs()
	ff.cmd = exec.Command("ffmpeg", args...)
	ff.cmd.Stdout = os.Stdout
	ff.cmd.Stderr = os.Stderr

	err = ff.cmd.Start()
	if err != nil {
		return nil, err
	}

	done = make(chan bool)

	go func(cmd *exec.Cmd) {
		cmd.Wait()
		ff.cmd = nil
		done <- true
	}(ff.cmd)

	return done, nil
}

func (ff *FFMPEG) isAlive() bool {
	if ff.file == "" {
		return false
	}
	fileInfo, err := os.Stat(ff.file)
	if err != nil {
		return false
	}
	return time.Since(fileInfo.ModTime()) < time.Second*15
}

func (ff *FFMPEG) Run(ctx context.Context) (err error) {

	var execDone chan bool
	start := true
	delay := 0
	deadCount := 0
	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-execDone:
			fmt.Println("ffmpeg has stopped -> restarting soon...")
			start = true
			delay = 1
		case <-ticker.C:
		}

		if delay > 0 {
			delay = delay - 1
		}

		if start && delay == 0 {
			start = false
			deadCount = 0
			execDone, err = ff.exec(ctx)
			if err != nil {
				fmt.Println("ffmpeg failed to start:", err)
				start = true
				delay = 1
			}
			continue
		}

		if ff.cmd == nil || ff.isAlive() {
			deadCount = 0
			continue
		}

		deadCount = deadCount + 1

		if deadCount >= 3 {
			fmt.Println("ffmpeg is running, but stopped writing files -> restart now")
			ff.cmd.Process.Kill()
		} else {
			fmt.Println("ffmpeg is running, but stopped writing files -> restart soon")
		}
	}
}

func Run(ctx context.Context, dir, url string) {
	ff := FFMPEG{
		dir: dir,
		url: url,
	}
	ff.Run(ctx)
}
