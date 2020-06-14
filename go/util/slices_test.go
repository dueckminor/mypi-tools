package util

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestStringsContains(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(StringsContains([]string{"foo", "bar"}, "foo")).To(BeTrue())
}
