package tmx_test

import (
	"os"

	"image/png"

	. "github.com/manyminds/tmx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test public renderer", func() {
	Context("Test render", func() {
		validateMapWithImage := func(mapFile, imageFile string) {
			f, err := os.Open(mapFile)
			Expect(err).ToNot(HaveOccurred())
			testMap, err := NewMap(f)
			Expect(err).ToNot(HaveOccurred())
			c := NewImageCanvasFromMap(*testMap)
			renderer := NewRenderer(*testMap, c)
			err = renderer.Render(0)
			Expect(err).ToNot(HaveOccurred())
			f, err = os.Open(imageFile)
			Expect(err).ToNot(HaveOccurred())
			expected, err := png.Decode(f)
			Expect(err).ToNot(HaveOccurred())
			Expect(expected).To(EqualImage(c.Image()))
		}

		It("should render all layers when visible gzip", func() {
			validateMapWithImage("./testfiles/simple_example.tmx", "./testfiles/simple_example_expected.png")
		})

		It("should render only visible images zlib", func() {
			validateMapWithImage("./testfiles/simple_example_zlib.tmx", "./testfiles/simple_example_zlib_expected.png")
		})
	})
})
