package tmx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTmx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tmx Suite")
}
