package geometry

import (
	"errors"
)

type Tin struct {
	//Points     []Point
	//ConvexHull []Point
	Triangles []Triangle
	MaxX      float64
	MaxY      float64
	MinX      float64
	MinY      float64
	//Halfedges  []int
}

// Triangulate returns a Delaunay triangulation of the provided points.
func CreateTin(points []Point, nodata float64) (*Tin, error) {
	t := newTriangulator(points)
	var minx, miny, maxx, maxy float64
	minx = 180
	miny = 180
	maxx = -180
	maxy = -180
	err := t.triangulate()
	if err != nil {
		//return &Tin{points, t.convexHull(), t.triangles, t.halfedges}, err
		return &Tin{}, err
	}
	ts := t.triangles
	tris := make([]Triangle, 0)
	count := 0
	for i := 0; i < len(ts); i += 3 {
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		if p0.Z != nodata || p1.Z != nodata || p2.Z != nodata {
			if maxx < p0.X {
				maxx = p0.X
			}
			if maxx < p1.X {
				maxx = p1.X
			}
			if maxx < p2.X {
				maxx = p2.X
			}
			if minx > p0.X {
				minx = p0.X
			}
			if minx > p1.X {
				minx = p1.X
			}
			if minx > p2.X {
				minx = p2.X
			}
			if maxy < p0.Y {
				maxy = p0.Y
			}
			if maxy < p1.Y {
				maxy = p1.Y
			}
			if maxy < p2.Y {
				maxy = p2.Y
			}
			if miny > p0.Y {
				miny = p0.Y
			}
			if miny > p1.Y {
				miny = p1.Y
			}
			if miny > p2.Y {
				miny = p2.Y
			}
			tris = append(tris, Triangle{P1: p0, P2: p1, P3: p2})
			count++
		}
	}
	tris = tris[:count] //count-1?
	return &Tin{Triangles: tris, MaxX: maxx, MinX: minx, MaxY: maxy, MinY: miny}, err
}
func (t *Tin) ComputeValue(x float64, y float64) (float64, error) {
	for _, tri := range t.Triangles {
		v, err := tri.GetValue(x, y)
		if err == nil {
			return v, err
		}
	}
	return -9999, errors.New("Point was not in triangles.")
}

/*
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
*/
