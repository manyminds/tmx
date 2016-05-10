package tmx_test

import (
	"errors"
	"image"
	"io"
	"os"

	. "github.com/manyminds/tmx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nilResourceLoader struct {
}

func (m nilResourceLoader) LocateResource(filepath string) (image.Image, error) {
	return nil, nil
}

type moduloErrorResourceLoader struct {
	callCount int
}

var errAlreadyLoaded = errors.New("Already loaded")

func (m *moduloErrorResourceLoader) LocateResource(filepath string) (image.Image, error) {
	m.callCount++
	if m.callCount%2 == 0 {
		return nil, errAlreadyLoaded
	}

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

		It("will find a testfile", func() {
			loader := FilesystemLocator{}
			_, err := loader.LocateResource("testfiles/chipset.png")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("test lazy resource locator", func() {
		BeforeEach(func() {
			c, err := os.Open("testfiles/chipset.png")
			Expect(err).ToNot(HaveOccurred())
			t, err := os.Create("testfiles/chipset_copy.png")
			Expect(err).ToNot(HaveOccurred())
			_, err = io.Copy(t, c)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.Remove("testfiles/chipset_copy.png")
		})

		It("caches files even after remove", func() {
			loader := NewLazyResourceLocator(FilesystemLocator{})
			img, err := loader.LocateResource("testfiles/chipset_copy.png")
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
			err = os.Remove("testfiles/chipset_copy.png")
			Expect(err).ToNot(HaveOccurred())
			img, err = loader.LocateResource("testfiles/chipset_copy.png")
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())

			manager, ok := loader.(ResourceManager)
			Expect(ok).To(BeTrue())
			manager.UnsetResource("testfiles/chipset_copy.png")
			img, err = loader.LocateResource("testfiles/chipset_copy.png")
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})

		It("will only load a specific resource once", func() {
			loader := &moduloErrorResourceLoader{}
			lazyLoader := NewLazyResourceLocator(loader)
			_, err := lazyLoader.LocateResource("lol")
			Expect(err).ToNot(HaveOccurred())
			_, err = lazyLoader.LocateResource("lol")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
