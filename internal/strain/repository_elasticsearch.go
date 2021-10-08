package strain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/ganvoa/biopipe-tools/internal"
)

type StrainRepository interface {
	SaveAll(strains []Strain) error
}

type Accesion struct {
	SeqInsert           int    `json:"seq_insert"`
	SeqLibrary          string `json:"seq_library"`
	ExperimentAccession string `json:"experiment_accession"`
	Accession           string `json:"accession"`
	SeqPlatform         string `json:"seq_platform"`
}

type Strain struct {
	SecondarySampleAccession string     `json:"secondary_sample_accession"`
	Comment                  string     `json:"comment"`
	Ecor                     string     `json:"ecor"`
	IsDownloadable           bool       `json:"is_downloadable"`
	AntibioticResistance     string     `json:"antibiotic_resistance"`
	Strain                   string     `json:"strain"`
	Postcode                 string     `json:"postcode"`
	Id                       int        `json:"id"`
	Best_assembly            int        `json:"best_assembly"`
	City                     string     `json:"city"`
	CollectionDate           int        `json:"collection_date"`
	ExtraRowInfo             string     `json:"extra_row_info"`
	CollectionMonth          int        `json:"collection_month"`
	NotEditable              bool       `json:"not_editable"`
	Owner                    int        `json:"owner"`
	Continent                string     `json:"continent"`
	Admin1                   string     `json:"admin1"`
	SourceDetails            string     `json:"source_details"`
	Admin2                   string     `json:"admin2"`
	Latitude                 float32    `json:"latitude"`
	StudyAccession           string     `json:"study_accession"`
	SerologicalGroup         string     `json:"serological_group"`
	SourceNiche              string     `json:"source_niche"`
	Barcode                  string     `json:"barcode"`
	Serotype                 string     `json:"serotype"`
	Uberstrain               int        `json:"uberstrain"`
	SimpleDisease            string     `json:"simple_disease"`
	SecondaryStudyAccession  string     `json:"secondary_study_accession"`
	CollectionYear           int        `json:"collection_year"`
	PathNonpath              string     `json:"path_nonpath"`
	AssemblyStatus           string     `json:"assembly_status"`
	Accession                []Accesion `json:"Accession"`
	Created                  string     `json:"created"`
	Country                  string     `json:"country"`
	ReleaseDate              string     `json:"release_date"`
	Disease                  string     `json:"disease"`
	Longitude                float32    `json:"longitude"`
	SampleAccession          string     `json:"sample_accession"`
	SimplePathogenesis       string     `json:"simple_pathogenesis"`
	SourceType               string     `json:"source_type"`
	Contact                  string     `json:"contact"`
	Species                  string     `json:"species"`
	CollectionTime           string     `json:"collection_time"`
	Downloaded               bool       `json:"downloaded"`
	IntegronFinder           bool       `json:"integron_finder"`
}

type strainRepositoryElasticSearch struct {
	index  string
	client *elasticsearch.Client
	logger internal.Logger
}

func NewRepository(index string, client *elasticsearch.Client, logger internal.Logger) *strainRepositoryElasticSearch {
	repository := new(strainRepositoryElasticSearch)
	repository.index = index
	repository.client = client
	repository.logger = logger
	return repository
}

func (repo strainRepositoryElasticSearch) SaveAll(strains []Strain) error {
	var bodyBuf bytes.Buffer

	for _, strain := range strains {

		preString := fmt.Sprintf(`{"index":{"_id": "%v"}}`, strain.Id)
		preBytes := []byte(preString)
		bodyBuf.Write(preBytes)
		bodyBuf.WriteByte('\n')

		strain.IntegronFinder = false
		strain.Downloaded = false
		strainId := strconv.Itoa(strain.Id)
		repo.logger.Debugf("strainId %s", strainId)
		createString, err := json.Marshal(strain)
		if err != nil {
			return err
		}
		bodyBuf.Write(createString)
		bodyBuf.WriteByte('\n')
	}

	req := esapi.BulkRequest{
		Index: repo.index,
		Body:  &bodyBuf,
	}
	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
