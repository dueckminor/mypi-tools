package setup

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestCreateBootImage(t *testing.T) {
	g := NewGomegaWithT(t)

	err := CreateBootImage(nil)
	g.Expect(err).To(BeNil())
}
