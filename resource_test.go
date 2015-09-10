package tmx_test

import (
	"image"

	. "github.com/manyminds/tmx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nilResourceLoader struct {
}

func (m nilResourceLoader) LocateResource(filepath string) (image.Image, error) {
	return nil, nil
}

var _ = Describe("Test public api", func() {
	Context("check types for interface interface", func() {
		It("can be implemented", func() {
			var loader interface{}
			loader = nilResourceLoader{}
			_, ok := loader.(ResourceLocator)
			Expect(ok).To(Equal(true))
		})

		It("is implemented by default", func() {
			var loader interface{}
			loader = FilesystemLocator{}
			_, ok := loader.(ResourceLocator)
			Expect(ok).To(Equal(true))
		})

		It("will error on invalid resource", func() {
			loader := FilesystemLocator{}
			_, err := loader.LocateResource("/null")
			Expect(err).To(HaveOccurred())
		})
	})
})
