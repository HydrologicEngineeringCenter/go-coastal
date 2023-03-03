package resultswriters

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type WoodHoleResultsWriter struct {
	filepath         string
	frequencies      []float64
	results          map[string]woodHoleStructureResult
	Layer            *gdal.Layer
	ds               *gdal.DataSource
	currentFrequency int
	discountFactor   float64
}
type woodHoleStructureResult struct {
	Name             string
	x                float64
	y                float64
	OccType          string
	DamCat           string
	Depths           []float64
	Waves            []float64
	StructureDamages []float64
	ContentDamages   []float64
}

func InitwoodHoleResultsWriterFromFile(filepath string, frequencies []float64, discountFactor float64) (*WoodHoleResultsWriter, error) {
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
			w := fmt.Sprintf("%v_wave", s)
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
			fieldDefw := gdal.CreateFieldDefinition(w, gdal.FT_Real)
			defer fieldDefw.Destroy()
			newLayer.CreateField(fieldDefw, true)
		}
		fieldDefsead := gdal.CreateFieldDefinition("s_Real_EAD", gdal.FT_Real)
		defer fieldDefsead.Destroy()
		newLayer.CreateField(fieldDefsead, true)
		fieldDefcead := gdal.CreateFieldDefinition("c_Real_EAD", gdal.FT_Real)
		defer fieldDefcead.Destroy()
		newLayer.CreateField(fieldDefcead, true)
		fieldDefasead := gdal.CreateFieldDefinition("s_Adj_EAD", gdal.FT_Real)
		defer fieldDefasead.Destroy()
		newLayer.CreateField(fieldDefasead, true)
		fieldDefacead := gdal.CreateFieldDefinition("c_Adj_EAD", gdal.FT_Real)
		defer fieldDefacead.Destroy()
		newLayer.CreateField(fieldDefacead, true)
	}()
	newLayer.StartTransaction()
	return &WoodHoleResultsWriter{filepath: filepath, results: t, frequencies: frequencies, Layer: &newLayer, ds: &dsOut, discountFactor: discountFactor}, nil
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
		xi, err := r.Fetch("x")
		if err != nil {
			//painic?
		}
		x := xi.(float64)
		yi, err := r.Fetch("y")
		if err != nil {
			//painic?
		}
		y := yi.(float64)
		wsr = woodHoleStructureResult{
			Name:             name,
			OccType:          occtype,
			x:                x,
			y:                y,
			DamCat:           damcat,
			StructureDamages: make([]float64, len(srw.frequencies)),
			ContentDamages:   make([]float64, len(srw.frequencies)),
			Depths:           make([]float64, len(srw.frequencies)),
			Waves:            make([]float64, len(srw.frequencies)),
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
	d, err := r.Fetch("hazard")
	if err != nil {
		//painic?
	}
	ce := d.(hazards.CoastalEvent)
	whsr.Depths[whrw.currentFrequency] = ce.Depth()
	whsr.Waves[whrw.currentFrequency] = ce.WaveHeight()
}
func (srw *WoodHoleResultsWriter) Close() {
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
			w := fmt.Sprintf("%v_wave", s)
			widx := layerDef.FieldIndex(w)
			feature.SetFieldFloat64(widx, r.Waves[i])
			//fmt.Println(s)

		}
		//c_EAD, s_EAD
		cead := compute.ComputeSpecialEAD(r.ContentDamages, srw.frequencies)
		ceadidx := layerDef.FieldIndex("c_Real_EAD")
		feature.SetFieldFloat64(ceadidx, cead)

		sead := compute.ComputeSpecialEAD(r.StructureDamages, srw.frequencies)
		seadidx := layerDef.FieldIndex("s_Real_ EAD")
		feature.SetFieldFloat64(seadidx, sead)

		acead := cead * srw.discountFactor
		aceadidx := layerDef.FieldIndex("c_Adj_EAD")
		feature.SetFieldFloat64(aceadidx, acead)

		asead := sead * srw.discountFactor
		aseadidx := layerDef.FieldIndex("s_Adj_ EAD")
		feature.SetFieldFloat64(aseadidx, asead)

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
