package heatmap

import (
	"image"
	"math"
)

func initImageIntensityItem(size image.RGBA) imageIntensityItem {
	ii := imageIntensityItem{}
	bounds := size.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		row := []int{}
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			row = append(row, 0)
		}
		ii = append(ii, row)
	}
	return ii
}

func findLimits(points []DataPoint) limits {
	minx, miny := points[0].X(), points[0].Y()
	maxx, maxy := minx, miny

	for _, p := range points {
		minx = math.Min(p.X(), minx)
		miny = math.Min(p.Y(), miny)
		maxx = math.Max(p.X(), maxx)
		maxy = math.Max(p.Y(), maxy)
	}

	return limits{apoint{minx, miny, 0}, apoint{maxx, maxy, 0}}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getMinMaxValues(ii imageIntensityItem) (int, int) {
	minVal := 0
	maxVal := 0
	for _, row := range ii {
		for _, col := range row {
			minVal = min(col, minVal)
			maxVal = max(col, maxVal)
		}
	}
	return minVal, maxVal
}
