package csvstreamer

import (
	"os"
	"testing"
	"time"
)

var DataCsv = []string{"./testdata/tdata1.csv", "./testdata/tdata2.csv"}

//TestHandlerMock default initializer
func TestHandlerMock(t *testing.T) {
	t.Log("Init test")

	for _, csvpath := range DataCsv {
		done := make(chan struct{})
		defer close(done)

		csv := New(csvpath)
		result, errc := csv.Parse(done)
		csvto := csvpath + "-mock-" + time.Now().Format("2006-01-02_150405.999") + ".csv"
		t.Log("SAVE::TO", csvto)
		csv2 := New(csvto)
		idx := 0

		var batchrow [][]string
		for p := range result {
			idx++
			batchrow = append(batchrow, p.Row)
			if len(batchrow) > 1024 {
				csv2.Save(true, batchrow...)
				batchrow = [][]string{}
			}
		}
		if len(batchrow) > 0 {
			csv2.Save(true, batchrow...)
		}
		t.Log("csv-total:", idx)
		if idx <= 0 {
			t.Fatal("Oops! must parse at least 1 row")
		}
		// Check whether the something happened
		if err := <-errc; err != nil {
			t.Fatal("Oops!", err)
		}
		//clear ;-)
		_ = os.Remove(csvto)
	}
	t.Log("OK")
}

//TestHandlerSimple default initializer
func TestHandlerSimple(t *testing.T) {
	t.Log("Init test")

	for _, csvpath := range DataCsv {

		csv := New(csvpath)
		result, errc := csv.Simple()
		csvto := csvpath + "-simple-" + time.Now().Format("2006-01-02_150405.999") + ".csv"
		t.Log("SAVE::TO", csvto)
		csv2 := New(csvto)
		idx := 0
		var batchrow [][]string
		for _, p := range result {
			idx++
			batchrow = append(batchrow, p.Row)
			if len(batchrow) > 1024 {
				csv2.Save(true, batchrow...)
				batchrow = [][]string{}
			}
		}
		if len(batchrow) > 0 {
			csv2.Save(true, batchrow...)
		}
		t.Log("csv-total:", idx)
		if idx <= 0 {
			t.Fatal("Oops! must parse at least 1 row")
		}
		// Check whether the something happened
		if errc != nil {
			t.Fatal("Oops!", errc)
		}
		//clear ;-)
		_ = os.Remove(csvto)
	}
	t.Log("OK")
}

//BenchmarkMock
func BenchmarkMock(t *testing.B) {
	t.Log("Init start", t.N)
	for i := 1; i <= t.N; i++ {
		for _, csvpath := range DataCsv {
			done := make(chan struct{})
			defer close(done)
			csv := New(csvpath)
			result, errc := csv.Parse(done)
			idx := 0
			for p := range result {
				idx++
				_ = p.Row
			}
			t.Log("csv-total:", idx)
			if idx <= 0 {
				t.Fatal("Oops! must parse at least 1 row")
			}
			// Check whether the something happened
			if err := <-errc; err != nil {
				t.Fatal("Oops!", err)
			}
		}
	}
	t.Log("OK")
}
