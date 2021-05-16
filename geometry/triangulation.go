package geometry

import (
	"errors"

	"github.com/tidwall/rtree"
)

type Tin struct {
	Triangles []Triangle
	MaxX      float64
	MaxY      float64
	MinX      float64
	MinY      float64
	Tree      rtree.RTree
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
	var tr rtree.RTree
	for i := 0; i < len(ts); i += 3 {
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		if p0.Z != nodata || p1.Z != nodata || p2.Z != nodata {
			t := CreateTriangle(p0, p1, p2)
			e := t.Extent()
			if maxx < e.UpperRight.X {
				maxx = e.UpperRight.X
			}
			if minx > e.LowerLeft.X {
				minx = e.LowerLeft.X
			}
			if maxy < e.UpperRight.Y {
				maxy = e.UpperRight.Y
			}
			if miny > e.LowerLeft.Y {
				miny = e.LowerLeft.Y
			}
			tris = append(tris, t)
			tr.Insert(e.LowerLeft.ToXY(), e.UpperRight.ToXY(), t)
			count++
		}
	}
	tris = tris[:count] //count-1?
	return &Tin{Triangles: tris, MaxX: maxx, MinX: minx, MaxY: maxy, MinY: miny, Tree: tr}, err
}
func (t *Tin) ComputeValue(x float64, y float64) (float64, error) {
	var v float64
	nodata := -9999.0
	var err error
	v = nodata
	t.Tree.Search([2]float64{x, y}, [2]float64{x, y},
		func(min, max [2]float64, value interface{}) bool {
			tri, ok := value.(Triangle)
			if ok {
				v, err = tri.GetValue(x, y)
				if err == nil {
					return false
				} else {
					return true
				}
			}
			return true
		},
	)
	if v == nodata {
		return nodata, errors.New("Point was not in triangles.")
	}
	if err == nil {
		return v, err
	}
	return nodata, errors.New("Point was not in triangles.")
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
