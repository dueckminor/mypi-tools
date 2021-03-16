package downloads

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestAlpineDownloader(t *testing.T) {
	g := NewGomegaWithT(t)

	d := NewAlpineDownloader()
	fd := d.GetDownloaderForVersion("3.13.2", "aarch64")

	g.Expect(fd.MetaData.Name).To(Equal("alpine-rpi-3.13.2-aarch64.tar.gz"))
	g.Expect(fd.MetaData.Checksums["sha256"]).To(Equal("3dc14236dec90078c5b989db305be7cf0aff0995c8cdb006dcccf13b0ac92f97"))

	// fd.StartDownload()
	// g.Eventually(fd.DownloadCompleted(), time.Minute*5).Should(BeTrue())
	// g.Eventually(fd.DownloadVerified(), time.Minute*5).Should(BeTrue())
}
