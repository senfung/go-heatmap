package main

import (
	"image"
	"image/png"
	"math/rand"
	"os"

	"github.com/senfung/go-heatmap"
	"github.com/senfung/go-heatmap/schemes"
)

func main() {
	points := []heatmap.DataPoint{}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			points = append(points, heatmap.P(float64(i), float64(j), rand.Int()))
		}
	}

	// scheme, _ := schemes.FromImage("../schemes/color_scheme.png")
	scheme := schemes.Classic

	img := heatmap.Heatmap(image.Rect(0, 0, 1024, 1024),
		points, 150, 128, scheme)
	png.Encode(os.Stdout, img)
}
