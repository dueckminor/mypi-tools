package setup

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	. "github.com/onsi/gomega"
)

var (
	dirTest = ""
)

func init() {
	pc, filename, _, _ := runtime.Caller(0)
	dirTest = path.Join(path.Dir(filename), "../../gen/test")
	err := os.MkdirAll(dirTest, os.ModePerm)
	if err != nil {
		panic(err)
	}

	f := runtime.FuncForPC(pc).Name()
	fmt.Println(f)
}

func TestNewFileWriter(t *testing.T) {
	g := NewGomegaWithT(t)

	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc).Name()
	fmt.Println(f)

	fw, err := NewFileWriter(dirTest)
	g.Expect(fw, err).NotTo(BeNil())

	w, err := fw.CreateFile(FileInfo{Type: FileTypeFile, Name: "test.txt"})
	g.Expect(w, err).NotTo(BeNil())

	g.Expect(w.Write([]byte("test"))).To(Equal(4))
	g.Expect(w.Close()).NotTo(HaveOccurred())

	w, err = fw.CreateFile(FileInfo{Type: FileTypeFile, Name: "../test.txt"})
	g.Expect(w, err).Error().To(HaveOccurred())

	w, err = fw.CreateFile(FileInfo{Type: FileTypeDir, Name: "subdir"})
	g.Expect(w, err).To(BeNil())

	err = fw.Close()
	g.Expect(err).To(BeNil())
}
