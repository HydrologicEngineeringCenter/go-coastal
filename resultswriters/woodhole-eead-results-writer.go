package resultswriters

import (
	"errors"
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/dewberry/gdal"
)

type WoodHoleEEADResultsWriter struct {
	filepath      string
	results       map[string]woodHoleStructureEEADResult
	Layer         *gdal.Layer
	ds            *gdal.DataSource
	AnalysisYears []int
}
type woodHoleStructureEEADResult struct {
	Name           string
	x              float64
	y              float64
	OccType        string
	DamCat         string
	AnalysisYears  []int
	StructureEEADs []float64
	ContentEEADs   []float64
}

func InitwoodHoleEEADResultsWriterFromFile(filepath string, AnalysisYears []int) (*WoodHoleEEADResultsWriter, error) {
	//make the maps
	t := make(map[string]woodHoleStructureEEADResult, 1)
	//create the geopackage
	driverOut := gdal.OGRDriverByName("GPKG")
	dsOut, okOut := driverOut.Create(filepath, []string{})
	if !okOut {
		//error out?
		return nil, errors.New("could not create output")
	}
	//defer dsOut.Destroy() -> probably should destroy on close?
	//set spatial reference?
	sr := gdal.CreateSpatialReference("")
	sr.FromEPSG(4326)
	newLayer := dsOut.CreateLayer("results", sr, gdal.GT_Point, []string{"GEOMETRY_NAME=shape"}) //forcing point data type.  source type (using lyaer.type()) from postgis was a generic geometry

	func() {
		fieldDef := gdal.CreateFieldDefinition("objectid", gdal.FT_Integer)
		defer fieldDef.Destroy()
		newLayer.CreateField(fieldDef, true)
	}()
	func() {
		fieldDefName := gdal.CreateFieldDefinition(Name, gdal.FT_String)
		defer fieldDefName.Destroy()
		newLayer.CreateField(fieldDefName, true)
		fieldDefx := gdal.CreateFieldDefinition(X, gdal.FT_Real)
		defer fieldDefx.Destroy()
		newLayer.CreateField(fieldDefx, true)
		fieldDefy := gdal.CreateFieldDefinition(Y, gdal.FT_Real)
		defer fieldDefy.Destroy()
		newLayer.CreateField(fieldDefy, true)
		fieldDefOT := gdal.CreateFieldDefinition(OccType, gdal.FT_String)
		defer fieldDefOT.Destroy()
		newLayer.CreateField(fieldDefOT, true)
		fieldDefDC := gdal.CreateFieldDefinition(DamCat, gdal.FT_String)
		defer fieldDefDC.Destroy()
		newLayer.CreateField(fieldDefDC, true)
		for _, val := range AnalysisYears {
			sd := fmt.Sprintf("%v_%v", StrucutreEquivalentEAD, val) //s for structure c for content
			//create field.
			fieldDefsead := gdal.CreateFieldDefinition(sd, gdal.FT_Real)
			defer fieldDefsead.Destroy()
			newLayer.CreateField(fieldDefsead, true)
			cd := fmt.Sprintf("%v_%v", ContentEquivalentEAD, val)
			fieldDefcead := gdal.CreateFieldDefinition(cd, gdal.FT_Real)
			defer fieldDefcead.Destroy()
			newLayer.CreateField(fieldDefcead, true)
		}
		//toteeadc toteeads
		sd := fmt.Sprintf("%v_%v", StrucutreEquivalentEAD, "tot")
		fieldDefeeads := gdal.CreateFieldDefinition(sd, gdal.FT_Real)
		defer fieldDefeeads.Destroy()
		newLayer.CreateField(fieldDefeeads, true)

		cd := fmt.Sprintf("%v_%v", ContentEquivalentEAD, "tot")
		fieldDefeeadc := gdal.CreateFieldDefinition(cd, gdal.FT_Real)
		defer fieldDefeeadc.Destroy()
		newLayer.CreateField(fieldDefeeadc, true)
	}()
	newLayer.StartTransaction()
	return &WoodHoleEEADResultsWriter{filepath: filepath, results: t, Layer: &newLayer, ds: &dsOut}, nil
}
func (srw *WoodHoleEEADResultsWriter) Write(r consequences.Result) {
	n, err := r.Fetch(Name)
	if err != nil {
		//painic?
	}
	name := n.(string)
	wsr, ok := srw.results[name]
	if !ok {
		//create on first pass.
		dc, err := r.Fetch(DamCat)
		if err != nil {
			//painic?
		}
		damcat := dc.(string)
		ot, err := r.Fetch(OccType)
		if err != nil {
			//painic?
		}
		occtype := ot.(string)
		//grab x and y
		xi, err := r.Fetch(X)
		if err != nil {
			//painic?
		}
		x := xi.(float64)
		yi, err := r.Fetch(Y)
		if err != nil {
			//painic?
		}
		y := yi.(float64)
		wsr = woodHoleStructureEEADResult{
			Name:           name,
			x:              x,
			y:              y,
			OccType:        occtype,
			DamCat:         damcat,
			AnalysisYears:  []int{},
			StructureEEADs: []float64{},
			ContentEEADs:   []float64{},
		}
	}
	//get first damage year and first eead values.
	ceeadi, err := r.Fetch(ContentEquivalentEAD)
	if err != nil {
		//painic?
	}
	ceead := ceeadi.(float64)
	seeadi, err := r.Fetch(StrucutreEquivalentEAD)
	if err != nil {
		//painic?
	}
	seead := seeadi.(float64)
	ai, err := r.Fetch(AnalysisYear)
	if err != nil {
		//painic?
	}
	a := ai.(int)
	wsr.ContentEEADs = append(wsr.ContentEEADs, ceead)
	wsr.StructureEEADs = append(wsr.StructureEEADs, seead)
	wsr.AnalysisYears = append(wsr.AnalysisYears, a)

	srw.results[name] = wsr

}
func (srw *WoodHoleEEADResultsWriter) Close() {
	layerDef := srw.Layer.Definition()

	//defer feature.Destroy()
	pointIndex := 0
	//rows

	for _, r := range srw.results {
		feature := layerDef.Create()
		defer feature.Destroy()
		fidx := layerDef.FieldIndex("objectid")
		feature.SetFieldInteger(fidx, pointIndex)
		g := gdal.Create(gdal.GT_Point)
		g.SetPoint(0, r.x, r.y, 0)
		feature.SetGeometryDirectly(g)
		//name
		sidx := layerDef.FieldIndex(Name)
		feature.SetFieldString(sidx, r.Name)
		//x
		xidx := layerDef.FieldIndex(X)
		feature.SetFieldFloat64(xidx, r.x)
		//y
		yidx := layerDef.FieldIndex(Y)
		feature.SetFieldFloat64(yidx, r.y)
		//occtype
		oidx := layerDef.FieldIndex(OccType)
		feature.SetFieldString(oidx, r.OccType)
		//damcat
		dcidx := layerDef.FieldIndex(DamCat)
		feature.SetFieldString(dcidx, r.DamCat)
		//frequency based headers
		totSEEAD := 0.0
		totCEEAD := 0.0
		for i, val := range r.AnalysisYears {
			sd := fmt.Sprintf("%v_%v", StrucutreEquivalentEAD, val) //s for structure c for content
			sidx := layerDef.FieldIndex(sd)
			feature.SetFieldFloat64(sidx, r.StructureEEADs[i])
			totSEEAD += r.StructureEEADs[i]
			cd := fmt.Sprintf("%v_%v", ContentEquivalentEAD, val)
			cidx := layerDef.FieldIndex(cd)
			feature.SetFieldFloat64(cidx, r.ContentEEADs[i])
			totCEEAD += r.ContentEEADs[i]
		}
		ceadidx := layerDef.FieldIndex(fmt.Sprintf("%v_%v", ContentEquivalentEAD, "tot"))
		feature.SetFieldFloat64(ceadidx, totCEEAD)

		seadidx := layerDef.FieldIndex(fmt.Sprintf("%v_%v", StrucutreEquivalentEAD, "tot"))
		feature.SetFieldFloat64(seadidx, totSEEAD)
		err := srw.Layer.Create(feature)
		if err != nil {
			fmt.Println(err)
		}
		pointIndex++
	}
	err2 := srw.Layer.CommitTransaction()
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Printf("Closing, wrote %v features\n", pointIndex-1)
	srw.ds.Destroy()
}
