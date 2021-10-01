package resultswriters

import (
	"fmt"
	"io"
	"os"

	"github.com/USACE/go-consequences/consequences"
)

type summaryResultsWriter struct {
	filepath   string
	w          io.Writer
	grandTotal float64
	totals     map[string]float64
}

func InitSummaryResultsWriterFromFile(filepath string) (*summaryResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return &summaryResultsWriter{}, err
	}
	//make the maps
	t := make(map[string]float64, 1)
	return &summaryResultsWriter{filepath: filepath, w: w, totals: t}, nil
}
func InitSummaryResultsWriter(w io.Writer) *summaryResultsWriter {
	t := make(map[string]float64, 1)
	return &summaryResultsWriter{filepath: "not applicapble", w: w, totals: t}
}
func (srw *summaryResultsWriter) Write(r consequences.Result) {
	//hardcoding for structures to experiment and think it through.
	var cat = "damage category"
	var structDam = "structure EAD"
	var contDam = "content EAD"
	var totDam = 0.0
	var damcat = ""
	h := r.Headers
	for i, v := range h {
		if v == cat {
			//add data to the map from this index in results
			damcat = r.Result[i].(string)
		}
		if v == structDam {
			totDam += r.Result[i].(float64)
		}
		if v == contDam {
			totDam += r.Result[i].(float64)
		}
	}
	srw.grandTotal += totDam
	//update damcat totals.
	t, ok := srw.totals[damcat]
	if ok {
		t += totDam
		srw.totals[damcat] = t
	} else {
		srw.totals[damcat] = totDam
	}
}
func (srw *summaryResultsWriter) Close() {
	fmt.Fprintf(srw.w, "Grand Total is %v\n", srw.grandTotal)
	h := srw.totals
	for i, v := range h {
		fmt.Fprintf(srw.w, "Damages for %v were %v\n", i, v)
	}
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
