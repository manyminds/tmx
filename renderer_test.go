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
		var (
			testMap *Map
			c       Canvas
		)

		BeforeEach(func() {
			f, err := os.Open("./testfiles/simple_example.tmx")
			Expect(err).ToNot(HaveOccurred())
			testMap, err = NewMap(f)
			Expect(err).ToNot(HaveOccurred())
			c = NewImageCanvasFromMap(*testMap)
		})

		It("should render all layers", func() {
			renderer := NewRenderer(*testMap, c)
			err := renderer.Render()
			Expect(err).ToNot(HaveOccurred())
			f, err := os.Open("./testfiles/simple_example_expected.png")
			Expect(err).ToNot(HaveOccurred())
			expected, err := png.Decode(f)
			Expect(err).ToNot(HaveOccurred())
			ic, ok := c.(*ImageCanvas)
			Expect(ok).To(Equal(true), "invalid type")
			Expect(expected).To(EqualImage(ic.Image()))
		})
	})
})
