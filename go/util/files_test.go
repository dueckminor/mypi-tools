package util

import (
	"os"
	"path"
	"testing"

	. "github.com/onsi/gomega"
)

func TestFileExists(t *testing.T) {
	g := NewWithT(t)
	tmpDir, err := os.MkdirTemp("", "testdir")
	g.Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(tmpDir)

	tmpFile := path.Join(tmpDir, "testfile")

	err = os.WriteFile(tmpFile, []byte("testdata"), 0644)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(FileExists(tmpDir)).To(BeFalse())
	g.Expect(FileExists(tmpFile)).To(BeTrue())
	g.Expect(FileExists(tmpFile + ".bak")).To(BeFalse())
}

func TestFileIsSafe(t *testing.T) {
	g := NewWithT(t)
	g.Expect(FileIsSafe("foo.bar")).To(BeTrue())
	g.Expect(FileIsSafe("..foo.bar")).To(BeFalse())
	g.Expect(FileIsSafe("/foo.bar")).To(BeFalse())
	g.Expect(FileIsSafe("../foo.bar")).To(BeFalse())
}

func TestFileIsSafePath(t *testing.T) {
	g := NewWithT(t)
	g.Expect(FileIsSafePath("foo.bar")).To(BeTrue())
	g.Expect(FileIsSafePath("..foo.bar")).To(BeFalse())
	g.Expect(FileIsSafePath("/foo.bar")).To(BeTrue())
	g.Expect(FileIsSafePath("../foo.bar")).To(BeFalse())
}
