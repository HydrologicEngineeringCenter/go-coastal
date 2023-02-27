package resultswriters

import (
	"io"
	"os"

	"github.com/USACE/go-consequences/consequences"
)

type woodHoleResultsWriter struct {
	filepath         string
	w                io.Writer
	frequencies      []float64
	results          map[string]woodHoleStructureResult
	currentFrequency int
}
type woodHoleStructureResult struct {
	Name             string
	OccType          string
	DamCat           string
	StructVal        float64
	ContentVal       float64
	Frequencies      []float64
	StructureDamages []float64
	ContentDamages   []float64
	StructureEAD     float64
	ContentEAD       float64
}

func InitwoodHoleResultsWriterFromFile(filepath string, frequencies []float64) (*woodHoleResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return &woodHoleResultsWriter{}, err
	}
	//make the maps
	t := make(map[string]woodHoleStructureResult, 1)
	return &woodHoleResultsWriter{filepath: filepath, w: w, results: t, frequencies: frequencies}, nil
}
func (srw *woodHoleResultsWriter) updateFrequencyIndex(i int) {
	srw.currentFrequency = i
}
func (srw *woodHoleResultsWriter) Write(r consequences.Result) {
	n, err := r.Fetch("name")
	if err != nil {
		//painic?
	}
	name := n.(string)
	wsr, ok := srw.results[name]
	if !ok {
		dc, err := r.Fetch("damcat")
		if err != nil {
			//painic?
		}
		damcat := dc.(string)
		ot, err := r.Fetch("occtype")
		if err != nil {
			//painic?
		}
		occtype := ot.(string)
		sv, err := r.Fetch("structVal")
		if err != nil {
			//panic
		}
		structVal := sv.(float64)
		cv, err := r.Fetch("contVal")
		if err != nil {
			//panic
		}
		contVal := cv.(float64)
		wsr = woodHoleStructureResult{
			Name:       name,
			OccType:    occtype,
			DamCat:     damcat,
			StructVal:  structVal,
			ContentVal: contVal,
		}
	}
	wsr.updateDamageInfo(r)
	srw.results[name] = wsr

}
func (wsr *woodHoleStructureResult) updateDamageInfo(r consequences.Result) {
	//use current frequency to set the appropriate value in the frequencies index.

}
func (srw *woodHoleResultsWriter) Close() {
	//fmt.Fprintf(srw.w, "Grand Total is %v\n", srw.grandTotal)
	//h := srw.totals
	//for i, v := range h {
	//fmt.Fprintf(srw.w, "Damages for %v were %v\n", i, v)
	//}
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
