package tmx

import "strconv"

type hexcolor string

//Tiled uses this defaults when no background color
//was set
const (
	defaultRed   = 128
	defaultBlue  = 128
	defaultGreen = 128
	defaultAlpha = 255
)

//implement image.Color interface
func (c hexcolor) RGBA() (r, g, b, a uint32) {
	if c == "" {
		return defaultRed, defaultBlue, defaultGreen, defaultAlpha
	}

	data := []byte(string(c))
	if len(data) == 0 {
		return defaultRed, defaultBlue, defaultGreen, defaultAlpha
	}

	if data[0] == '#' {
		data = data[1:]
	}

	if len(data) != 6 {
		return defaultRed, defaultBlue, defaultGreen, defaultAlpha
	}

	rx, err := strconv.ParseInt(string(data[0:2]), 16, 0)
	if err != nil {
		return defaultRed, defaultBlue, defaultGreen, defaultAlpha
	}

	gx, err := strconv.ParseInt(string(data[2:4]), 16, 0)
	if err != nil {
		return defaultRed, defaultBlue, defaultGreen, defaultAlpha
	}

	bx, err := strconv.ParseInt(string(data[4:6]), 16, 0)
	if err != nil {
		return defaultRed, defaultBlue, defaultGreen, defaultAlpha
	}

	return uint32(rx), uint32(gx), uint32(bx), 255
}
