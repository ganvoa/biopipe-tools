package integron

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

type IntegronParser struct {
	logger internal.Logger
}

func NewParser(logger internal.Logger) IntegronParser {
	parser := IntegronParser{}
	parser.logger = logger
	return parser
}

func (ip IntegronParser) Parse(path string) ([]string, error) {

	integrons := []string{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(path, file.Name())

		integronsFound, err := ip.fileParse(filePath)
		if err != nil {
			return nil, err
		}

		integronsChain, err := ip.toChain(integronsFound)
		if err != nil {
			return nil, err
		}

		integrons = append(integrons, integronsChain)

	}

	ip.logger.Debugf("integrons found %d", len(integrons))

	return integrons, nil
}

func (ip IntegronParser) toChain(integrons []Integron) (string, error) {
	delimiter := " "
	chain := ""

	for i := len(integrons) - 1; i >= 0; i-- {
		if i > 0 {
			chain = chain + integrons[i].Annotation + delimiter
		} else {
			chain = chain + integrons[i].Annotation
		}

	}

	return chain, nil
}

func (ip IntegronParser) fileParse(filePath string) ([]Integron, error) {
	ip.logger.Debugf("opening file %s", filePath)

	integrons := []Integron{}

	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Comment = '#'

	headers, err := r.Read()

	if err != nil {
		return nil, err
	}

	if len(headers) != 14 {
		return nil, errors.New("num of header cols must be 14")
	}

	for {

		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(record) != 14 {
			return nil, errors.New("num of record cols must be 14")
		}

		posBegin, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		posEnd, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, err
		}

		strand, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, err
		}

		ir := Integron{}
		ir.ID_Integron = record[0]
		ir.ID_Replicon = record[1]
		ir.Element = record[2]
		ir.Pos_Beg = posBegin
		ir.Pos_End = posEnd
		ir.Strand = strand
		ir.Evalue = record[6]
		ir.Type_Elt = record[7]
		ir.Annotation = record[8]
		ir.Model = record[9]
		ir.Type = record[10]
		ir.Default = record[11]
		ir.Distance_2attC = record[12]
		ir.Considered_Topology = record[13]

		integrons = append(integrons, ir)
	}

	return integrons, nil
}
