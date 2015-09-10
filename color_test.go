package tmx_test

import (
	. "github.com/manyminds/tmx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test Background Color parsing", func() {
	Context("Test default cases", func() {
		It("Should be default if empty", func() {
			m := Map{}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(128)))
			Expect(g).To(Equal(uint32(128)))
			Expect(b).To(Equal(uint32(128)))
			Expect(a).To(Equal(uint32(255)))
		})

		It("Should be default if too long", func() {
			m := Map{BackgroundColor: "thisisnocolor"}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(128)))
			Expect(g).To(Equal(uint32(128)))
			Expect(b).To(Equal(uint32(128)))
			Expect(a).To(Equal(uint32(255)))
		})

		It("Should be default if invalid", func() {
			m := Map{BackgroundColor: "nocolor"}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(128)))
			Expect(g).To(Equal(uint32(128)))
			Expect(b).To(Equal(uint32(128)))
			Expect(a).To(Equal(uint32(255)))
		})

		It("will default if it fails to decode red", func() {
			m := Map{BackgroundColor: "fgaaaa"}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(128)))
			Expect(g).To(Equal(uint32(128)))
			Expect(b).To(Equal(uint32(128)))
			Expect(a).To(Equal(uint32(255)))
		})

		It("will default if it fails to decode green", func() {
			m := Map{BackgroundColor: "faagaa"}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(128)))
			Expect(g).To(Equal(uint32(128)))
			Expect(b).To(Equal(uint32(128)))
			Expect(a).To(Equal(uint32(255)))
		})

		It("will default if it fails to decode blue", func() {
			m := Map{BackgroundColor: "faaaga"}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(128)))
			Expect(g).To(Equal(uint32(128)))
			Expect(b).To(Equal(uint32(128)))
			Expect(a).To(Equal(uint32(255)))
		})

		It("Should be ok if correct", func() {
			m := Map{BackgroundColor: "#23af40"}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(35)))
			Expect(g).To(Equal(uint32(175)))
			Expect(b).To(Equal(uint32(64)))
			Expect(a).To(Equal(uint32(255)))
		})

		It("Should be ok if hash is missing", func() {
			m := Map{BackgroundColor: "23AF40"}

			r, g, b, a := m.BackgroundColor.RGBA()
			Expect(r).To(Equal(uint32(35)))
			Expect(g).To(Equal(uint32(175)))
			Expect(b).To(Equal(uint32(64)))
			Expect(a).To(Equal(uint32(255)))
		})
	})
})
