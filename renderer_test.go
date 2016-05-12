package tmx_test

import (
	"fmt"
	"os"

	"image/png"

	. "github.com/manyminds/tmx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test public renderer", func() {
	Context("Test render", func() {
		validateMapWithImage := func(mapFile, imageFile string, elapsedTime int64) {
			f, err := os.Open(mapFile)
			Expect(err).ToNot(HaveOccurred())
			testMap, err := NewMap(f)
			Expect(err).ToNot(HaveOccurred())
			c := NewImageCanvasFromMap(*testMap)
			renderer := NewRenderer(*testMap, c)
			err = renderer.Render(elapsedTime)
			Expect(err).ToNot(HaveOccurred())
			//			this code generates reference images if necessary
			/*
			 *
			 *      target, err := os.Create(imageFile)
			 *      if err != nil {
			 *        target, err = os.Open(imageFile)
			 *        if err != nil {
			 *          log.Fatal(err)
			 *        }
			 *      }
			 *      defer target.Close()
			 *
			 *      err = png.Encode(target, c.Image())
			 *      if err != nil {
			 *        log.Fatal(err)
			 *      }
			 *
			 */
			f, err = os.Open(imageFile)
			Expect(err).ToNot(HaveOccurred())
			expected, err := png.Decode(f)
			Expect(err).ToNot(HaveOccurred())
			Expect(expected).To(EqualImage(c.Image()))
		}

		It("should render all layers when visible gzip", func() {
			validateMapWithImage("./testfiles/simple_example.tmx", "./testfiles/simple_example_expected.png", 0)
		})

		It("should render only visible images zlib", func() {
			validateMapWithImage("./testfiles/simple_example_zlib.tmx", "./testfiles/simple_example_zlib_expected.png", 0)
		})

		It("should render non squared uncompressed maps", func() {
			validateMapWithImage("./testfiles/uncompressed_not_square.tmx", "./testfiles/uncompressed_not_square.png", 0)
		})

		It("renders animated tiles", func() {
			validateMapWithImage("./testfiles/animated_example_zlib.tmx", "./testfiles/animated_example_zlib_01.png", 0)
			validateMapWithImage("./testfiles/animated_example_zlib.tmx", "./testfiles/animated_example_zlib_02.png", 101)
			validateMapWithImage("./testfiles/animated_example_zlib.tmx", "./testfiles/animated_example_zlib_03.png", 201)
			validateMapWithImage("./testfiles/animated_example_zlib.tmx", "./testfiles/animated_example_zlib_04.png", 301)
		})
	})

	Context("Test flip mode", func() {
		It("will have a working String", func() {
			Expect(fmt.Sprintf("%s", FlipNone)).To(Equal("None"))
			Expect(fmt.Sprintf("%s", FlipHorizontal)).To(Equal("Horizontal"))
			Expect(fmt.Sprintf("%s", FlipVertical)).To(Equal("Vertical"))
			Expect(fmt.Sprintf("%s", FlipDiagonal)).To(Equal("Diagonal"))
		})
	})
})
