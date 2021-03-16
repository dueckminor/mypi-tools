package downloads

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

type FileMetadata struct {
	Name      string            `json:"name"`
	URL       string            `json:"url"`
	Size      int64             `json:"size"`
	Checksums map[string]string `json:"checksums,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
}

type FileDownloader struct {
	MetaData      *FileMetadata
	mutex         sync.RWMutex
	targetFile    string
	receivedBytes int64
	verified      bool
}

type Downloader struct {
	mutex           sync.RWMutex
	TargetDir       string
	FileMetadatas   []*FileMetadata
	FileDownloaders map[string]*FileDownloader
}

func downloadFromURL(url string, w io.Writer) (err error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadToFile(url string, filepath string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	return downloadFromURL(url, out)
}

func downloadToBuffer(url string) (data []byte, err error) {
	var b bytes.Buffer
	err = downloadFromURL(url, bufio.NewWriter(&b))
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func downloadToString(url string) (data string, err error) {
	b, err := downloadToBuffer(url)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func downloadGetSize(url string) (downloadSize int64, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Got wrong status: %d", resp.StatusCode)
	}
	return resp.ContentLength, nil
}

func (d *Downloader) GetMetadatas() []*FileMetadata {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.FileMetadatas
}

func (d *Downloader) GetMetadata(filename string) (m *FileMetadata) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, m = range d.FileMetadatas {
		if m.Name == filename {
			return m
		}
	}
	return nil
}

func (d *Downloader) GetDownloader(filename string) (fd *FileDownloader) {
	m := d.GetMetadata(filename)
	if m == nil {
		return nil
	}

	d.mutex.RLock()
	fd = d.FileDownloaders[filename]
	d.mutex.RUnlock()
	if fd != nil {
		return fd
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()

	fd = &FileDownloader{
		MetaData:   m,
		targetFile: path.Join(d.TargetDir, m.Name),
	}
	if nil == d.FileDownloaders {
		d.FileDownloaders = make(map[string]*FileDownloader)
	}
	d.FileDownloaders[filename] = fd
	return fd
}

func (d *Downloader) StartDownload(filename string) (fd *FileDownloader) {
	fd = d.GetDownloader(filename)
	if fd == nil {
		return nil
	}
	fd.StartDownload()
	return fd
}

func (fd *FileDownloader) GetTargetFile() string {
	return fd.targetFile
}

func (fd *FileDownloader) StartDownload() {
	if fd.DownloadCompleted() {
		return
	}
	go func() {
		fd.mutex.Lock()
		defer fd.mutex.Unlock()
		if fd.DownloadCompleted() {
			return
		}
		err := downloadToFile(fd.MetaData.URL, fd.targetFile)
		if err != nil {
			os.Remove(fd.targetFile)
			return
		}
	}()
}

func (fd *FileDownloader) DownloadCompleted() bool {
	stat, err := os.Stat(fd.targetFile)
	if err != nil {
		return false
	}
	if fd.MetaData.Size == stat.Size() {
		return true
	}
	return false
}

func (fd *FileDownloader) DownloadVerified() bool {
	if fd.verified {
		return true
	}
	if !fd.DownloadCompleted() {
		return false
	}

	f, err := os.Open(fd.targetFile)
	if err != nil {
		return false
	}
	defer f.Close()

	var hasher hash.Hash

	for alg, sum := range fd.MetaData.Checksums {
		switch alg {
		case "sha256":
			hasher = sha256.New()
		}
		if hasher == nil {
			continue
		}

		_, err = io.Copy(hasher, f)
		if err != nil {
			return false
		}

		if sum == hex.EncodeToString(hasher.Sum(nil)) {
			fd.verified = true
			break
		}
	}

	return fd.verified
}
