package csvstreamer

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	//VersionMajor main ver no.
	VersionMajor = "0.1"
	//VersionMinor sub  ver no.
	VersionMinor = "0"
	//DecimalMaxLen float max len when converted to string
	DecimalMaxLen = 20
	//EnclosedBy default enclosing-chars of the csv cols
	EnclosedBy = "\""
	//DefaultEmpty is the default string if empty
	DefaultEmpty = "-"
	//DefaultComma is the default separator of the elements for writing
	DefaultComma = ","
)

var (
	//BuildTime pass during build time
	BuildTime string
	//UtilsVersion is the app ver string
	UtilsVersion string
)

//internal system initialize
func init() {
	//uniqueness
	rand.Seed(time.Now().UnixNano())
	if BuildTime == "" {
		BuildTime = time.Now().Format("20060102150405.999")
	}
	UtilsVersion = "Ver: " + VersionMajor + "." + VersionMinor + "-" + BuildTime
}

//CsvStream data mapping
type CsvStream struct {
	filename string
	emptyval string
	enclosed string
	sep      string
}

//CsvResult the stream of csv data
type CsvResult struct {
	Row []string
}

//CsvStreamer the stream map
type CsvStreamer interface {
	Version() string
	Format(rows ...interface{}) []string
	ToCsvStr(msg []string) string
	Parse(done chan struct{}) (chan CsvResult, chan error)
	Simple() ([]CsvResult, error)
	Save(data ...[]string) error
	Append(data ...[]string) error
}

//New the initializer
func New(csv string, setters ...CsvStreamOpt) *CsvStream {
	//defaults
	me := &CsvStream{
		filename: strings.TrimSpace(csv),
		emptyval: DefaultEmpty,
		enclosed: EnclosedBy,
		sep:      DefaultComma,
	}
	//maybe params are passed
	for _, setter := range setters {
		setter(me)
	}
	//good ;-)
	return me
}

//Version return the utils version
func (c *CsvStream) Version() string {
	return UtilsVersion
}

//Save the list of items into new CSV file
func (c *CsvStream) Save(data ...[]string) error {
	fh, err := os.OpenFile(c.filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	//sanity
	defer fh.Close()
	w := csv.NewWriter(fh)
	for _, record := range data {
		if err := w.Write(record); err != nil {
			return err
		}
	}
	// Write any buffered data to the underlying writer
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

//Append the list of items into old CSV file or new if not exists
func (c *CsvStream) Append(data ...[]string) error {
	fh, err := os.OpenFile(c.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	//sanity
	defer fh.Close()
	w := csv.NewWriter(fh)
	for _, record := range data {
		if err := w.Write(record); err != nil {
			return err
		}
	}
	// Write any buffered data to the underlying writer
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

//Parse handler on parsing the CSV
func (c *CsvStream) Parse(done chan struct{}) (chan CsvResult, chan error) {
	res := make(chan CsvResult)
	ret := make(chan error, 1)
	go func() {
		var wg sync.WaitGroup
		errc := func() error {
			csvFile, err := os.Open(c.filename)
			if err != nil {
				return err
			}
			reader := csv.NewReader(bufio.NewReader(csvFile))
			for {
				line, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					return error
				}
				wg.Add(1)
				go func() { // HL
					select {
					case res <- CsvResult{Row: line}:
					case <-done: // HL
					}
					wg.Done()
				}()

			}
			//cancel if needed
			select {
			case <-done: // HL
				return errors.New("parse canceled")
			default:
				return nil
			}
		}()
		//goroutine to close result once all the sends are done.
		go func() { // HL
			wg.Wait()
			close(res)
		}()
		// No select needed here, since errc is buffered.
		ret <- errc // HL
	}()

	return res, ret
}

//Simple handler on parsing the CSV no channel
func (c *CsvStream) Simple() ([]CsvResult, error) {
	var res []CsvResult
	csvFile, err := os.Open(c.filename)
	if err != nil {
		return res, err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return res, error
		}

		res = append(res, CsvResult{Row: line})

	}
	return res, nil
}

//Format format all fields
func (c *CsvStream) Format(rows ...interface{}) []string {
	//defaults
	msg := []string{}
	for _, x := range rows {
		switch x.(type) {
		case string:
			msg = append(msg, c.enclosed+fmtDoubleQts(strings.TrimSpace(fmt.Sprintf("%s", x)))+c.enclosed)
		case int, int32, int64:
			msg = append(msg, c.enclosed+fmt.Sprintf("%d", x)+"\"")
		case float64, float32:
			msg = append(msg, c.enclosed+fmtNumeric(fmt.Sprintf("%f", x))+c.enclosed)
		default: //unknown field
			msg = append(msg, c.enclosed+c.emptyval+c.enclosed)
		}
	}
	return msg
}

//ToCsvStr convert the row list into 1 csv string
func (c *CsvStream) ToCsvStr(msg []string) string {
	//defaults
	return strings.Join(msg, c.sep)
}
