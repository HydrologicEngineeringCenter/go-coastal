package geometry

import "math"

type Point struct {
	X         float64
	Y         float64
	HasZValue bool
	Z         float64 //make a slice?
}

func (a Point) squaredDistance(b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func (a Point) distance(b Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

func (a Point) sub(b Point) Point {
	return Point{X: a.X - b.X, Y: a.Y - b.Y}
}
func (a Point) ToXY() [2]float64 {
	return [2]float64{a.X, a.Y}
}
