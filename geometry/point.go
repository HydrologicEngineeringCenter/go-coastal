package geometry

import "math"

type PointZ struct {
	*Point
	Z []float64 //make a slice?
}
type PointZZ struct {
	*Point
	ZSwl  []float64 //make a slice?
	ZHm0  []float64 //make a slice?
	ZElev float64
}
type Point struct {
	X float64
	Y float64
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
func (a PointZ) squaredDistance(b PointZ) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func (a PointZ) distance(b PointZ) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

func (a PointZ) sub(b PointZ) Point {
	return Point{X: a.X - b.X, Y: a.Y - b.Y}
}
func (a PointZ) ToXY() [2]float64 {
	return [2]float64{a.X, a.Y}
}
func (a PointZ) ToPoint() [2]float64 {
	return [2]float64{a.X, a.Y}
}

func (a PointZZ) squaredDistance(b PointZZ) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func (a PointZZ) distance(b PointZZ) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

func (a PointZZ) sub(b PointZZ) Point {
	return Point{X: a.X - b.X, Y: a.Y - b.Y}
}
func (a PointZZ) ToXY() [2]float64 {
	return [2]float64{a.X, a.Y}
}
func (a PointZZ) ToPoint() [2]float64 {
	return [2]float64{a.X, a.Y}
}
