package resultswriters

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/dewberry/gdal"
)

type GridWriter struct {
	filepath string
	ds       *gdal.Dataset
	band     *gdal.RasterBand
}

func InitGridWriterFromFIle(filepath string, xSteps int, ySteps int, xmin float64, xdelta float64, ymin float64, ydelta float64) (*GridWriter, error) {
	//make the maps
	driver, _ := gdal.GetDriverByName("GTiff")
	outdata := driver.Create(filepath, xSteps, ySteps, 1, gdal.Float32, []string{"BIGTIFF=YES", "TILED=YES", "COMPRESS=DEFLATE", "TILESIZE=256", "ZLEVEL=1"}) // TILED=YES COMPRESS=DEFLATE TILESIZE=256 ZLEVEL=1
	srs := gdal.CreateSpatialReference("")
	_ = srs.FromEPSG(4326)
	proj, _ := srs.ToWKT()
	outdata.SetProjection(proj)
	outdata.SetGeoTransform([6]float64{xmin, xdelta, 0, ymin, 0, ydelta})
	outband := outdata.RasterBand(1) //? 0 or 1?
	//value := 2.23//needs to be unsafe pointer
	//outband.WriteBlock(1,2,value)//need to use inverse geo transform to set the index for x and y.
	gw := GridWriter{filepath: filepath, band: &outband, ds: &outdata}
	return &gw, nil
}

func (srw *GridWriter) Write(r consequences.Result) {
	d, err := r.Fetch("depth")
	if err != nil {
		//painic?
	}
	depth := d.(float64)
	//need x and y locations
	xval, err := r.Fetch("x")
	if err != nil {
		//painic?
	}
	yval, err := r.Fetch("y")
	if err != nil {
		//painic?
	}
	x := xval.(float64)
	y := yval.(float64)
	igt := srw.ds.InvGeoTransform()
	px := int(igt[0] + x*igt[1] + y*igt[2])
	py := int(igt[3] + x*igt[4] + y*igt[5])
	buffer := make([]float64, 1*1)
	buffer[0] = depth
	srw.band.IO(gdal.Write, px, py, 1, 1, buffer, 1, 1, 0, 0)
	//srw.band.WriteBlock(px, py, unsafe.SliceData(buffer))//need 1.20 go or above.

}
func (srw *GridWriter) Close() {
	fmt.Printf("Closing")
	srw.band.GetDataset().Close()
}
