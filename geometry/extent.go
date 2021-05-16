package geometry

type Extent struct {
	LowerLeft  Point
	UpperRight Point
}

func (e Extent) Min() [2]float64 {
	return e.LowerLeft.ToXY()
}
func (e Extent) Max() [2]float64 {
	return e.UpperRight.ToXY()
}
