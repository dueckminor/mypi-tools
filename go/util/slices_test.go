package util

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestStringsContains(t *testing.T) {
	g := NewWithT(t)
	g.Expect(StringsContains([]string{"foo", "bar"}, "foo")).To(BeTrue())
	g.Expect(StringsContains([]string{"foo", "bar"}, "foobar")).To(BeFalse())
}

func TestStringsContainsAll(t *testing.T) {
	g := NewWithT(t)
	g.Expect(StringsContainsAll([]string{"foo", "bar"}, []string{"bar", "foo"})).To(BeTrue())
	g.Expect(StringsContainsAll([]string{"foo", "bar"}, []string{"bar", "FOO"})).To(BeFalse())
}
