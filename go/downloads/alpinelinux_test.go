package downloads

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestAlpineDownloader(t *testing.T) {
	g := NewGomegaWithT(t)

	d := NewAlpineDownloader()
	fd := d.GetDownloaderForVersion("3.13.2", "aarch64")

	g.Expect(fd.MetaData.Name).To(Equal("alpine-rpi-3.13.2-aarch64.tar.gz"))

	fd.StartDownload()

	g.Eventually(fd.DownloadCompleted(), time.Minute*5).Should(BeTrue())
	g.Eventually(fd.DownloadVerified(), time.Minute*5).Should(BeTrue())

}
