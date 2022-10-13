package analyzer

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/disintegration/imaging"
)

func loadImage(path string) (string, error) {
	img, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}
	img = imaging.Resize(img, 16, 16, imaging.NearestNeighbor)

	sbuilder := strings.Builder{}
	nrgba := img.(*image.NRGBA)
	for x := 0; x < nrgba.Bounds().Dx(); x++ {
		for y := 0; y < nrgba.Bounds().Dy(); y++ {
			sbuilder.WriteString(colorToHex(nrgba.At(x, y)))
			if y < nrgba.Bounds().Dy()-1 || x < nrgba.Bounds().Dx()-1 {
				sbuilder.WriteString(",")
			}
		}
	}

	return sbuilder.String(), nil
}

func colorToHex(color color.Color) string {
	r, g, b, _ := color.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", int(float64(r)/65535*255), int(float64(g)/65535*255), int(float64(b)/65535*255))
}
