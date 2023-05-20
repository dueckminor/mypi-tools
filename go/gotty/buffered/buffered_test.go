package buffered

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

var (
	lorem_text = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam
nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam`
	// voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita
	// kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem
	// ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
	// tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At
	// vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd
	// gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet`

	lorem_bytes = []byte(lorem_text)
)

func runTest(t *testing.T, write_size int, reader_count int, read_size int) {
	g := NewGomegaWithT(t)

	b, _ := NewBufferedTty()

	done := make(chan []byte)

	lorem_text = "000_001_002_003_004_005_006_007_008_009\n"
	lorem_text += "010_011_012_013_014_015_016_017_018_019\n"
	lorem_text += "020_021_022_023_024_025_026_027_028_029\n"
	lorem_text += "030_031_032_033_034_035_036_037_038_039\n"
	lorem_text += "040_041_042_043_044_045_046_047_048_049\n"
	lorem_bytes = []byte(lorem_text)

	end := func(have, size int) int {
		if have+size > len(lorem_bytes) {
			return len(lorem_bytes)
		}
		return have + size
	}

	reader := func() {
		s, _ := b.GetFactory().New(nil)
		read_buffer := make([]byte, len(lorem_bytes))
		read_count := 0
		for read_count < len(lorem_bytes) {
			now, _ := s.Read(read_buffer[read_count:end(read_count, read_size)])
			if now == 0 {
				break
			}
			read_count += now
			time.Sleep(time.Millisecond)
		}
		done <- read_buffer[:read_count]
	}

	for i := 0; i < reader_count; i++ {
		go reader()
	}

	write_count := 0
	for write_count < len(lorem_bytes) {
		now, _ := b.Write(lorem_bytes[write_count:end(write_count, write_size)])
		write_count += now
		g.Expect(now > 0).To(BeTrue())
		time.Sleep(time.Millisecond)
	}
	g.Expect(write_count).To(Equal(len(lorem_bytes)))

	for i := 0; i < reader_count; i++ {
		read_buffer := <-done
		g.Expect(read_buffer).To(Equal(lorem_bytes))
	}
}

func TestBuffered(t *testing.T) {
	runTest(t, 24, 10, 20)
}
