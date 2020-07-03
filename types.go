package heatmap

// A DataPoint to be plotted.
// These are all normalized to use the maximum amount of
// space available in the output image.
type DataPoint interface {
	X() float64
	Y() float64
	Value() int
}

type apoint struct {
	x     float64
	y     float64
	value int
}

func (a apoint) X() float64 {
	return a.x
}

func (a apoint) Y() float64 {
	return a.y
}

func (a apoint) Value() int {
	return a.value
}

// P is a shorthand simple datapoint constructor.
func P(x, y float64, value int) DataPoint {
	return apoint{x, y, value}
}

type imageIntensityItem [][]int

type limits struct {
	Min DataPoint
	Max DataPoint
}

func (l limits) inRange(lx, hx, ly, hy float64) bool {
	return l.Min.X() >= lx &&
		l.Max.X() <= hx &&
		l.Min.Y() >= ly &&
		l.Max.Y() <= hy
}

func (l limits) Dx() float64 {
	return l.Max.X() - l.Min.X()
}

func (l limits) Dy() float64 {
	return l.Max.Y() - l.Min.Y()
}
