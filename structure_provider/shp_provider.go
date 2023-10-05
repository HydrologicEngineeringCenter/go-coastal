package structure_provider

import (
	"math"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type uniformPointDataSet struct {
	xdelta float64
	ydelta float64
}

type depthReceptor struct {
	x float64
	y float64
}

func (s depthReceptor) Location() geography.Location {
	return geography.Location{X: s.x, Y: s.y}
}
func (dr depthReceptor) Compute(h hazards.HazardEvent) (consequences.Result, error) {
	headers := []string{"x", "y", "depth"}
	values := []interface{}{dr.x, dr.y, 0.0}
	if h.Depth() >= 0 {
		values = []interface{}{dr.x, dr.y, h.Depth()}
	}
	result := consequences.Result{
		Headers: headers,
		Result:  values,
	}
	return result, nil
}

func InitUniformPointDS(xdelta float64, ydelta float64) (uniformPointDataSet, error) {
	ds := uniformPointDataSet{
		xdelta: xdelta,
		ydelta: ydelta,
	}
	return ds, nil
}

// ByFips a streaming service for structure stochastic based on a bounding box
func (shp uniformPointDataSet) ByFips(fipscode string, sp consequences.StreamProcessor) {
	return
}

// ByBbox allows a shapefile to be streamed by bounding box
func (shp uniformPointDataSet) ByBbox(bbox geography.BBox, sp consequences.StreamProcessor) {
	shp.processBboxStream(bbox, sp)
}
func (shp uniformPointDataSet) processBboxStream(bbox geography.BBox, sp consequences.StreamProcessor) {
	y := 0
	x := 0
	//get distance in x domain
	xdist := bbox.Bbox[0] - bbox.Bbox[2]

	//get distance in y domain
	ydist := bbox.Bbox[1] - bbox.Bbox[3]
	//get total number of x and y steps
	xSteps := int(math.Floor(math.Abs(xdist) / shp.xdelta))
	ySteps := int(math.Floor(math.Abs(ydist) / shp.ydelta))
	//offset by half in each direction
	currentYval := bbox.Bbox[3] + (shp.ydelta / 2)
	var currentXval float64
	//generate a full row, incriment y and start the next row.
	for y < ySteps { //iterate across all rows
		x = 0
		currentXval = bbox.Bbox[0] + (shp.xdelta / 2)
		for x < xSteps { // Iterate across all x values in a row

			r := depthReceptor{
				x: currentXval,
				y: currentYval,
			}
			sp(r)
			x++
			currentXval += shp.xdelta
		}
		y++ //step to next row
		currentYval += shp.ydelta
	}
}
