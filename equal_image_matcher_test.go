package tmx_test

import (
	"fmt"
	"image"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

//EqualImage matches an image with an actual image
func EqualImage(expected interface{}) types.GomegaMatcher {
	return &EqualImageMatcher{
		Expected: expected,
	}
}

//EqualImageMatcher matches one image against another
type EqualImageMatcher struct {
	Expected interface{}
}

//Match matches Expected with actual
//types must be image.Image otherwise it fails
func (matcher *EqualImageMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}

	a, ok := actual.(image.Image)
	if !ok {
		return false, fmt.Errorf("actual must be of type image.Image")
	}

	e, ok := matcher.Expected.(image.Image)
	if !ok {
		return false, fmt.Errorf("expected must be of type image.Image")
	}

	if a.Bounds().Dx() != e.Bounds().Dx() {
		return false, fmt.Errorf("Different width values %d != %d", a.Bounds().Dx(), e.Bounds().Dx())
	}

	if a.Bounds().Dy() != e.Bounds().Dy() {
		return false, fmt.Errorf("Different height values %d != %d", a.Bounds().Dy(), e.Bounds().Dy())
	}

	for x := 0; x < a.Bounds().Dx(); x++ {
		for y := 0; y < a.Bounds().Dy(); y++ {
			if a.At(x, y) != e.At(x, y) {
				return false, fmt.Errorf("pixels don't match at %d %d. %d != %d", x, y, a.At(x, y), e.At(x, y))
			}
		}
	}

	return true, nil
}

//FailureMessage if images won't match
func (matcher *EqualImageMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to equal", matcher.Expected)
}

//NegatedFailureMessage if images match if they shouldn't
func (matcher *EqualImageMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}
