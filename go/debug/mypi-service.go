package debug

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/dueckminor/mypi-tools/go/util/network"
)

type MypiService struct {
	name    string
	ctxGo   context.Context
	cmdGo   *exec.Cmd
	ctxNode context.Context
	cmdNode *exec.Cmd
}

func NewMypiService(name string) (s *MypiService) {
	return &MypiService{name: name}
}

func (s *MypiService) StartNode() (err error) {
	s.ctxNode = context.Background()
	s.cmdNode = exec.CommandContext(s.ctxNode, "npm", "install")
	s.cmdNode.Stdout = os.Stdout
	s.cmdNode.Stderr = os.Stderr
	s.cmdNode.Dir = path.Join(GetWorkspaceRoot(), "web", s.name)
	s.cmdNode.Run()
	s.cmdNode.Wait()

	err = s.ctxNode.Err()
	if err != nil {
		return err
	}

	port, err := network.GetFreePort()
	if err != nil {
		return err
	}

	s.cmdNode = exec.CommandContext(s.ctxNode, path.Join(GetWorkspaceRoot(), "web", s.name, "node_modules/.bin/vue-cli-service"), "serve", "--port", strconv.Itoa(port))
	s.cmdNode.Stdout = os.Stdout
	s.cmdNode.Stderr = os.Stderr
	s.cmdNode.Dir = path.Join(GetWorkspaceRoot(), "web", s.name)
	err = s.cmdNode.Run()
	err = s.cmdNode.Wait()

	return nil
}

func getGoFiles(dir string) (filenames []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasSuffix(name, ".go") {
				filenames = append(filenames, path.Join(dir, file.Name()))
			}
		}
	}
	return filenames, nil
}

func (s *MypiService) StartGo() (err error) {
	s.ctxGo = context.Background()

	args := []string{"run"}
	goFiles, err := getGoFiles(path.Join(GetWorkspaceRoot(), "cmd", s.name))
	if err != nil {
		return err
	}

	port, err := network.GetFreePort()
	if err != nil {
		return err
	}

	args = append(args, goFiles...)

	args = append(args, fmt.Sprintf("--port=%v", port))
	args = append(args, fmt.Sprintf("--webpack-debug=http://localhost:%v", port+1))
	args = append(args, "--mypi-root=/opt/mypi")

	s.cmdNode = exec.CommandContext(s.ctxGo, "go", args...)
	s.cmdNode.Stdout = os.Stdout
	s.cmdNode.Stderr = os.Stderr
	s.cmdNode.Dir = GetWorkspaceRoot()
	s.cmdNode.Run()
	s.cmdNode.Wait()

	return nil
}
