package geometry

import (
	"errors"
	"fmt"

	"github.com/tidwall/rtree"
)

type Tin struct {
	MaxX float64
	MaxY float64
	MinX float64
	MinY float64
	Tree rtree.RTree
	Hull Polygon
	zidx int
}

// Triangulate returns a Delaunay triangulation of the provided points.
func CreateTin(points []PointZ, nodata float64, hull Polygon, zidx int) (*Tin, error) {
	t := newTriangulator(points)
	var minx, miny, maxx, maxy float64
	minx = 180
	miny = 180
	maxx = -180
	maxy = -180
	fmt.Println("Triangulating...")
	err := t.triangulate()
	if err != nil {
		return &Tin{}, err
	}
	ts := t.triangles
	count := 0
	var tr rtree.RTree
	for i := 0; i < len(ts); i += 3 {
		p0 := &points[ts[i+0]]
		p1 := &points[ts[i+1]]
		p2 := &points[ts[i+2]]
		lenz := len(p0.Z) - 1
		if p0.Z[lenz] != nodata || p1.Z[lenz] != nodata || p2.Z[lenz] != nodata {
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
			tr.Insert(e.LowerLeft.ToXY(), e.UpperRight.ToXY(), t)
			count++
		}
	}
	fmt.Println(fmt.Sprintf("Found %v triangles.", count))
	return &Tin{MaxX: maxx, MinX: minx, MaxY: maxy, MinY: miny, Tree: tr, Hull: hull}, err
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
				v, err = tri.GetValue(x, y, t.zidx)
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
