package tmx_test

import (
	. "github.com/manyminds/tmx"
	"github.com/manyminds/tmx/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests for map", func() {

	Context("Test map utility methods", func() {
		var testMap Map
		var miniTileset spec.Tileset
		var giantTileset spec.Tileset

		BeforeEach(func() {
			miniImage := spec.Image{Width: 320, Height: 320}
			miniTileset = spec.Tileset{
				FirstGID:   1,
				TileWidth:  32,
				TileHeight: 32,
				Image:      miniImage,
			}

			giantImage := spec.Image{Width: 16000, Height: 16000}
			giantTileset = spec.Tileset{
				FirstGID:   spec.GID(miniTileset.GetNumTiles()),
				TileWidth:  32,
				TileHeight: 32,
				Image:      giantImage,
			}

			testMap = Map{
				Tilesets:   []spec.Tileset{miniTileset, giantTileset},
				TileHeight: 32,
				TileWidth:  32,
				Width:      10,
				Height:     10,
			}
		})

		It("Should be a valid map", func() {
			err := testMap.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return no tileset for 0", func() {
			set, err := testMap.GetTilesetForGID(0)
			Expect(set).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should error for invalid gid range", func() {
			set, err := testMap.GetTilesetForGID(250100)
			Expect(err).To(HaveOccurred())
			Expect(set).To(BeNil())
		})

		It("should be returning miniTileset", func() {
			set, err := testMap.GetTilesetForGID(2)
			Expect(err).ToNot(HaveOccurred())
			Expect(set).To(Equal(&miniTileset))
		})

		It("should be returning miniTileset lower boundary", func() {
			set, err := testMap.GetTilesetForGID(100)
			Expect(err).ToNot(HaveOccurred())
			Expect(set).To(Equal(&miniTileset))
		})

		It("should be returning giantTileset", func() {
			set, err := testMap.GetTilesetForGID(101)
			Expect(err).ToNot(HaveOccurred())
			Expect(set).To(Equal(&giantTileset))
		})

		It("should be returning giantTileset upper boundary", func() {
			set, err := testMap.GetTilesetForGID(250099)
			Expect(err).ToNot(HaveOccurred())
			Expect(set).To(Equal(&giantTileset))
		})
	})
})
