package integron

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/ganvoa/biopipe-tools/internal"
)

type Integron struct {
	ID_Integron         string
	ID_Replicon         string
	Element             string
	Pos_Beg             int
	Pos_End             int
	Strand              int
	Evalue              string
	Type_Elt            string
	Annotation          string
	Model               string
	Type                string
	Default             string
	Distance_2attC      string
	Considered_Topology string
}

type integronParser struct {
	path      string
	integrons []Integron
	persist   Persister
	logger    internal.Logger
}

func NewParser(path string, persist Persister, logger internal.Logger) integronParser {
	parser := integronParser{}
	parser.path = path
	parser.logger = logger
	parser.persist = persist
	parser.integrons = []Integron{}
	return parser
}

func (ip integronParser) Parse() error {

	ip.logger.Debugf("opening file %s", ip.path)

	file, err := os.Open(ip.path)

	if err != nil {
		return err
	}

	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Comment = '#'

	headers, err := r.Read()

	if err != nil {
		return err
	}

	if len(headers) != 14 {
		return errors.New("num of header cols must be 14")
	}

	for {

		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if len(record) != 14 {
			return errors.New("num of record cols must be 14")
		}

		posBegin, err := strconv.Atoi(record[3])
		if err != nil {
			return err
		}

		posEnd, err := strconv.Atoi(record[4])
		if err != nil {
			return err
		}

		Strand, err := strconv.Atoi(record[5])
		if err != nil {
			return err
		}

		ir := Integron{}
		ir.ID_Integron = record[0]
		ir.ID_Replicon = record[1]
		ir.Element = record[2]
		ir.Pos_Beg = posBegin
		ir.Pos_End = posEnd
		ir.Strand = Strand
		ir.Evalue = record[6]
		ir.Type_Elt = record[7]
		ir.Annotation = record[8]
		ir.Model = record[9]
		ir.Type = record[10]
		ir.Default = record[11]
		ir.Distance_2attC = record[12]
		ir.Considered_Topology = record[13]
		ip.integrons = append(ip.integrons, ir)
	}

	ip.logger.Infof("integrons found %d", len(ip.integrons))

	return nil
}

func (ip integronParser) Save() error {
	err := ip.persist.Save(ip.integrons)
	if err != nil {
		return err
	}
	return nil
}
