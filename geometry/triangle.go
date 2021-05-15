package geometry

import "errors"

type Triangle struct {
	P1   Point
	P2   Point
	P3   Point
	HasZ bool
}

//https://codeplea.com/triangular-interpolation
func (t Triangle) GetValue(x float64, y float64) (float64, error) {
	invDenom := 1 / ((t.P2.Y-t.P3.Y)*(t.P1.X-t.P3.X) + (t.P3.X-t.P2.X)*(t.P1.Y-t.P3.Y))
	w1 := ((t.P2.Y-t.P3.Y)*(x-t.P3.X) + (t.P3.X-t.P2.X)*(y-t.P3.Y)) * invDenom
	w2 := ((t.P3.Y-t.P1.Y)*(x-t.P3.X) + (t.P1.X-t.P3.X)*(y-t.P3.Y)) * invDenom
	w3 := 1.0 - w1 - w2
	if w1 >= 0 && w2 >= 0 && w3 >= 0 {
		return (w1*t.P1.Z + w2*t.P2.Z + w3*t.P3.Z), nil
	}
	return -9999, errors.New("Point Outside Triangle")
}
