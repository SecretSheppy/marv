package colour

import (
	"strconv"
	"strings"
)

func Brightness(hex string) (float64, error) {
	hex = strings.TrimPrefix(hex, "#")
	r, err := strconv.ParseInt(hex[0:2], 16, 32)
	if err != nil {
		return 0, err
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 32)
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 32)
	if err != nil {
		return 0, err
	}
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b), nil
}

func IsBright(hex string) (bool, error) {
	b, err := Brightness(hex)
	if err != nil {
		return false, err
	}
	return b >= 128, nil
}
