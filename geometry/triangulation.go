package geometry

/*
type Tin struct {
	Triangles []Triangle
	//extent?
	HasZ bool
}
type Extent struct {
	LowerLeft  Point
	LowerRight Point
	UpperLeft  Point
	UpperRight Point
}
func (e Extent) Width() float64{
	return e.LowerRight.X - e.LowerLeft.X
}
func (e Extent) Height() float64{
	return e.LowerRight.Y - e.UpperRight.Y
}
func DelauneyTriangulation(points []Point, minx float64, maxx float64, miny float64, maxy float64) Tin {
	pointcount := len(points)
	//
	ext := Extent{
		UpperLeft: Point{X:minx,Y:maxy},
		UpperRight: Point{X:maxx,Y:maxy},
		LowerLeft: Point{X:minx,Y:miny},
		LowerRight: Point{X:maxx,Y:miny},
	}
	points = append(points, AddOctagonPointsAndTriangles(ext)...)
	MassPointInsertationStartingIndex := len(points)
	return Tin{}
}
func AddOctagonPointsAndTriangles(ext Extent)[]Point{
	// add in the points?
	w := ext.Width()
	h := ext.Height()
	delta := h/.5
	if w>h {
		delta = w/.5
	}
	points := make([]Point, 8)
	p0 := Point{X: (ext.UpperLeft.X-ext.UpperRight.X)/2, Y: ext.UpperLeft.Y + delta*2}
	p1 := Point{X: ext.UpperLeft.X + delta*2,Y: ext.LowerRight.Y + delta*2}
	p2 := Point{X: ext.UpperLeft.X+ delta*2, Y:(ext.UpperLeft.Y - ext.LowerLeft.Y)/2}
	p3 := Point{X: ext.UpperLeft.X + delta*2,Y:ext.UpperRight.Y - delta*2}

	p4 := Point{X: (ext.UpperLeft.X-ext.UpperRight.X)/2, Y: ext.LowerLeft.Y - delta*2}
	p5 := Point{X: ext.UpperRight.X - delta*2,Y: ext.UpperRight.Y - delta*2}
	p6 := Point{X: ext.UpperRight.X - delta*2, Y:(ext.UpperLeft.Y - ext.LowerLeft.Y)/2}
	p7 := Point{X: ext.UpperRight.X - delta*2,Y:ext.UpperRight.Y - delta*2}
	points[0] = p0
	points[1] = p1
	points[2] = p2
	points[3] = p3
	points[4] = p4
	points[5] = p5
	points[6] = p6
	points[7] = p7
	return points
}
func DelauneyInsertPoint(p Point){

}
*/

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Triangulation struct {
	Points     []Point
	ConvexHull []Point
	Triangles  []int
	Halfedges  []int
}

// Triangulate returns a Delaunay triangulation of the provided points.
func Triangulate(points []Point) (*Triangulation, error) {
	t := newTriangulator(points)
	err := t.triangulate()
	return &Triangulation{points, t.convexHull(), t.triangles, t.halfedges}, err
}

func (t *Triangulation) area() float64 {
	var result float64
	points := t.Points
	ts := t.Triangles
	for i := 0; i < len(ts); i += 3 {
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		result += area(p0, p1, p2)
	}
	return result / 2
}
func (t *Triangulation) Json(outpath string) {
	w, err := os.OpenFile(outpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	defer w.Close()
	if err != nil {
		panic(err)
	}
	s := "{\"type\": \"FeatureCollection\",\"features\": ["
	points := t.Points
	ts := t.Triangles
	for i := 0; i < len(ts); i += 3 {
		fmt.Fprintf(w, s)
		s = ""
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		s = "{\"type\": \"Feature\",\"geometry\": {\"type\": \"Polygon\",\"coordinates\": [["
		s += "[" + fmt.Sprintf("%g, %g", p0.X, p0.Y) + "],"
		s += "[" + fmt.Sprintf("%g, %g", p1.X, p1.Y) + "],"
		s += "[" + fmt.Sprintf("%g, %g", p2.X, p2.Y) + "],"
		s += "[" + fmt.Sprintf("%g, %g", p0.X, p0.Y) + "]]"
		s += "]},\"properties\": {\"p\": \"p\"}},"
	}
	s = strings.TrimRight(s, ",")
	s += "]}"
	fmt.Fprintf(w, s)
}

// Validate performs several sanity checks on the Triangulation to check for
// potential errors. Returns nil if no issues were found. You normally
// shouldn't need to call this function but it can be useful for debugging.
func (t *Triangulation) Validate() error {
	// verify halfedges
	for i1, i2 := range t.Halfedges {
		if i1 != -1 && t.Halfedges[i1] != i2 {
			return fmt.Errorf("invalid halfedge connection")
		}
		if i2 != -1 && t.Halfedges[i2] != i1 {
			return fmt.Errorf("invalid halfedge connection")
		}
	}

	// verify convex hull area vs sum of triangle areas
	hull1 := t.ConvexHull
	hull2 := ConvexHull(t.Points)
	area1 := polygonArea(hull1)
	area2 := polygonArea(hull2)
	area3 := t.area()
	if math.Abs(area1-area2) > 1e-9 || math.Abs(area1-area3) > 1e-9 {
		return fmt.Errorf("hull areas disagree: %f, %f, %f", area1, area2, area3)
	}

	// verify convex hull perimeter
	perimeter1 := polygonPerimeter(hull1)
	perimeter2 := polygonPerimeter(hull2)
	if math.Abs(perimeter1-perimeter2) > 1e-9 {
		return fmt.Errorf("hull perimeters disagree: %f, %f", perimeter1, perimeter2)
	}

	return nil
}
