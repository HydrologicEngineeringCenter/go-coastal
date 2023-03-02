package resultswriters

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/dewberry/gdal"
)

type WoodHoleResultsWriter struct {
	filepath         string
	frequencies      []float64
	results          map[string]woodHoleStructureResult
	Layer            *gdal.Layer
	ds               *gdal.DataSource
	currentFrequency int
}
type woodHoleStructureResult struct {
	Name             string
	x                float64
	y                float64
	OccType          string
	DamCat           string
	Depths           []float64
	StructureDamages []float64
	ContentDamages   []float64
}

func InitwoodHoleResultsWriterFromFile(filepath string, frequencies []float64) (*WoodHoleResultsWriter, error) {
	//make the maps
	t := make(map[string]woodHoleStructureResult, 1)
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
		fieldDefName := gdal.CreateFieldDefinition("name", gdal.FT_String)
		defer fieldDefName.Destroy()
		newLayer.CreateField(fieldDefName, true)
		fieldDefx := gdal.CreateFieldDefinition("x", gdal.FT_Real)
		defer fieldDefx.Destroy()
		newLayer.CreateField(fieldDefx, true)
		fieldDefy := gdal.CreateFieldDefinition("y", gdal.FT_Real)
		defer fieldDefy.Destroy()
		newLayer.CreateField(fieldDefy, true)
		fieldDefOT := gdal.CreateFieldDefinition("occtype", gdal.FT_String)
		defer fieldDefOT.Destroy()
		newLayer.CreateField(fieldDefOT, true)
		fieldDefDC := gdal.CreateFieldDefinition("damcat", gdal.FT_String)
		defer fieldDefDC.Destroy()
		newLayer.CreateField(fieldDefDC, true)
		//headers
		for _, val := range frequencies {
			s := strconv.FormatFloat(val, 'f', 3, 64)
			s = strings.Replace(s, "0.", ".", 1)
			sd := fmt.Sprintf("%v_%v_dam", s, "s") //s for structure c for content
			cd := fmt.Sprintf("%v_%v_dam", s, "c")
			d := fmt.Sprintf("%v_depth", s)
			//fmt.Println(s)
			fieldDefsd := gdal.CreateFieldDefinition(sd, gdal.FT_Real)
			defer fieldDefsd.Destroy()
			newLayer.CreateField(fieldDefsd, true)
			fieldDefcd := gdal.CreateFieldDefinition(cd, gdal.FT_Real)
			defer fieldDefcd.Destroy()
			newLayer.CreateField(fieldDefcd, true)
			fieldDefd := gdal.CreateFieldDefinition(d, gdal.FT_Real)
			defer fieldDefd.Destroy()
			newLayer.CreateField(fieldDefd, true)
		}
		fieldDefsead := gdal.CreateFieldDefinition("s_EAD", gdal.FT_String)
		defer fieldDefsead.Destroy()
		newLayer.CreateField(fieldDefsead, true)
		fieldDefcead := gdal.CreateFieldDefinition("c_EAD", gdal.FT_String)
		defer fieldDefcead.Destroy()
		newLayer.CreateField(fieldDefcead, true)
	}()
	newLayer.StartTransaction()
	return &WoodHoleResultsWriter{filepath: filepath, results: t, frequencies: frequencies, Layer: &newLayer, ds: &dsOut}, nil
}
func (srw *WoodHoleResultsWriter) UpdateFrequencyIndex(i int) {
	srw.currentFrequency = i
}
func (srw *WoodHoleResultsWriter) Write(r consequences.Result) {
	n, err := r.Fetch("fd_id")
	if err != nil {
		//painic?
	}
	name := n.(string)
	wsr, ok := srw.results[name]
	if !ok {
		//create on first pass.
		dc, err := r.Fetch("damage category")
		if err != nil {
			//painic?
		}
		damcat := dc.(string)
		ot, err := r.Fetch("occupancy type")
		if err != nil {
			//painic?
		}
		occtype := ot.(string)
		//grab x and y
		wsr = woodHoleStructureResult{
			Name:             name,
			OccType:          occtype,
			x:                0.0,
			y:                0.0,
			DamCat:           damcat,
			StructureDamages: make([]float64, len(srw.frequencies)),
			ContentDamages:   make([]float64, len(srw.frequencies)),
		}
	}
	wsr.updateDamageInfo(r, srw)
	srw.results[name] = wsr

}
func (whsr *woodHoleStructureResult) updateDamageInfo(r consequences.Result, whrw *WoodHoleResultsWriter) {
	//use current frequency to set the appropriate value in the frequencies index.
	sd, err := r.Fetch("structure damage")
	if err != nil {
		//painic?
	}
	sdam := sd.(float64)
	whsr.StructureDamages[whrw.currentFrequency] = sdam
	cd, err := r.Fetch("content damage")
	if err != nil {
		//painic?
	}
	cdam := cd.(float64)
	whsr.ContentDamages[whrw.currentFrequency] = cdam
}
func (srw *WoodHoleResultsWriter) Close() {
	layerDef := srw.Layer.Definition()
	feature := layerDef.Create()
	//defer feature.Destroy()
	pointIndex := 0
	//rows
	for _, r := range srw.results {
		feature.SetFieldInteger(0, pointIndex)
		g := gdal.Create(gdal.GT_Point)
		g.SetPoint(0, r.x, r.y, 0)
		feature.SetGeometryDirectly(g)
		//name
		sidx := layerDef.FieldIndex("name")
		feature.SetFieldString(sidx, r.Name)
		//x
		xidx := layerDef.FieldIndex("x")
		feature.SetFieldFloat64(xidx, r.x)
		//y
		yidx := layerDef.FieldIndex("y")
		feature.SetFieldFloat64(yidx, r.y)
		//occtype
		oidx := layerDef.FieldIndex("occtype")
		feature.SetFieldString(oidx, r.OccType)
		//damcat
		dcidx := layerDef.FieldIndex("damcat")
		feature.SetFieldString(dcidx, r.DamCat)
		//frequency based headers
		for i, val := range srw.frequencies {
			s := strconv.FormatFloat(val, 'f', 3, 64)
			s = strings.Replace(s, "0.", ".", 1)
			sd := fmt.Sprintf("%v_%v_dam", s, "s") //s for structure c for content
			sidx := layerDef.FieldIndex(sd)
			feature.SetFieldFloat64(sidx, r.StructureDamages[i])
			cd := fmt.Sprintf("%v_%v_dam", s, "c")
			cidx := layerDef.FieldIndex(cd)
			feature.SetFieldFloat64(cidx, r.ContentDamages[i])
			d := fmt.Sprintf("%v_depth", s)
			didx := layerDef.FieldIndex(d)
			feature.SetFieldFloat64(didx, r.Depths[i])
			//fmt.Println(s)

		}
		//c_EAD, s_EAD
		cead := compute.ComputeSpecialEAD(r.ContentDamages, srw.frequencies)
		ceadidx := layerDef.FieldIndex("c_EAD")
		feature.SetFieldFloat64(ceadidx, cead)

		sead := compute.ComputeSpecialEAD(r.StructureDamages, srw.frequencies)
		seadidx := layerDef.FieldIndex("s_EAD")
		feature.SetFieldFloat64(seadidx, sead)
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
