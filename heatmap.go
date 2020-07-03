// Package heatmap generates heatmaps for map overlays.
package heatmap

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
)

// Heatmap draws a heatmap.
//
// size is the size of the image to crate
// dotSize is the impact size of each point on the output
// opacity is the alpha value (0-255) of the impact of the image overlay
// scheme is the color palette to choose from the overlay
func Heatmap(size image.Rectangle, points []DataPoint, dotSize int, opacity uint8,
	scheme []color.Color) image.Image {

	dot := mkDot(float64(dotSize))
	limits := findLimits(points)

	// Draw black/alpha into the image
	bw := image.NewRGBA(size)
	intensityGraph := initImageIntensityItem(*bw)
	intensityGraph = placePoints(size, limits, bw, points, dot, intensityGraph)
	_, maxVal := getMinMaxValues(intensityGraph)
	rv := image.NewRGBA(size)

	// Then we transplant the pixels one at a time pulling from our color map
	warm(rv, opacity, scheme, intensityGraph, maxVal)
	return rv
}

func placePoints(size image.Rectangle, limits limits, bw *image.RGBA, points []DataPoint, dot draw.Image, ii imageIntensityItem) imageIntensityItem {
	for _, p := range points {
		ii = limits.placePoint(p, bw, dot, ii, size)
	}
	return ii
}

func warm(out draw.Image, opacity uint8, colors []color.Color, intensityGraph imageIntensityItem, maxVal int) {
	draw.Draw(out, out.Bounds(), image.Transparent, image.ZP, draw.Src)
	bounds := out.Bounds()
	collen := float64(len(colors)) // 751
	wg := &sync.WaitGroup{}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				colNum := len(intensityGraph[x])
				val := intensityGraph[x][colNum-y-1]
				// _, _, _, alpha := col.RGBA()
				if val > 0 {
					percent := float64(val) / float64(maxVal)
					template := colors[int((collen-1)*(1.0-percent))]
					tr, tg, tb, ta := template.RGBA()
					ta /= 256
					outalpha := uint8(float64(ta) * (float64(opacity) * math.Sqrt(percent) / 256.0))
					outcol := color.NRGBA{
						uint8(tr / 256),
						uint8(tg / 256),
						uint8(tb / 256),
						uint8(outalpha)}
					out.Set(x, y, outcol)
				}
			}
		}(x)
	}
	wg.Wait()
}

func mkDot(size float64) draw.Image {
	i := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
	md := 0.5 * math.Sqrt(math.Pow(float64(size)/2.0, 2)+math.Pow((float64(size)/2.0), 2))
	for x := float64(0); x < size; x++ {
		for y := float64(0); y < size; y++ {
			d := math.Sqrt(math.Pow(x-size/2.0, 2) + math.Pow(y-size/2.0, 2))
			if d < md {
				rgbVal := uint8(200.0*d/md + 50.0)
				rgba := color.NRGBA{0, 0, 0, 255 - rgbVal}

				// plsWork := uint8((255 - int(200.0*d/md+50.0)) / scale)
				// plsWork := uint8(255 / scale)
				// rgba := color.NRGBA{0, 0, 0, plsWork}
				i.Set(int(x), int(y), rgba)
			}
		}
	}

	return i
}

func (l limits) translate(p DataPoint, i draw.Image, dotsize int) (rv image.Point) {
	// Normalize to 0-1
	x := float64(p.X()-l.Min.X()) / float64(l.Dx())
	y := float64(p.Y()-l.Min.Y()) / float64(l.Dy())

	// And remap to the image
	rv.X = int(x * float64((i.Bounds().Max.X - dotsize)))
	rv.Y = int((1.0 - y) * float64((i.Bounds().Max.Y - dotsize)))

	return
}

func (l limits) placePoint(p DataPoint, i, dot draw.Image, ii imageIntensityItem, size image.Rectangle) imageIntensityItem {
	pos := l.translate(p, i, dot.Bounds().Max.X)
	dotw, doth := dot.Bounds().Max.X, dot.Bounds().Max.Y
	newImage := image.NewRGBA(size)
	draw.Draw(newImage, image.Rect(pos.X, pos.Y, pos.X+dotw, pos.Y+doth), dot, image.ZP, draw.Over)
	bounds := newImage.Bounds()
	wg := &sync.WaitGroup{}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				col := newImage.At(x, y)
				_, _, _, alpha := col.RGBA()
				if alpha > 0 {
					ii[x][y] = ii[x][y] + int(alpha)*p.Value()
				}
			}
		}(x)
	}
	wg.Wait()
	return ii
}
