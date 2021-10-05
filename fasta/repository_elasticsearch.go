package fasta

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type DownloadedRequest struct {
	Source struct {
		Downloaded bool `json:"downloaded"`
	} `json:"source"`
}

type Strain struct {
	Id         int `json:"id"`
	AssemblyId int `json:"best_assembly"`
}

type NotDownloadedResponse struct {
	Hits struct {
		Hits []struct {
			Source Strain `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type FastaRepositoryElasticSearch struct {
	index  string
	client *elasticsearch.Client
}

func NewRepository(index string, client *elasticsearch.Client) *FastaRepositoryElasticSearch {
	repository := new(FastaRepositoryElasticSearch)
	repository.index = index
	repository.client = client
	return repository
}

func (repo FastaRepositoryElasticSearch) NotDownloaded() ([]Strain, error) {
	reader := strings.NewReader(`{
		"query": {
		  "bool": {
			"must_not": [
			  {
				"term": {
				  "collection_year": {
					"value": 0
				  }
				}
			  },
			  {
				"term": {
				  "source_type.keyword": {
					"value": ""
				  }
				}
			  },
			  {
				"term": {
				  "is_downloadable": {
					"value": false
				  }
				}
			  },
			  {
				"term": {
				  "downloaded": {
					"value": true
				  }
				}
			  }
			]
		  }
		},
		"_source": [
		  "id",
		  "best_assembly",
		  "is_downloadable",
		  "collection_year",
		  "downloaded"
		]
	  }`)

	res, err := repo.client.Search(
		repo.client.Search.WithContext(context.Background()),
		repo.client.Search.WithIndex(repo.index),
		repo.client.Search.WithBody(reader),
		repo.client.Search.WithTrackTotalHits(true),
		repo.client.Search.WithFilterPath("hits.hits._source"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := NotDownloadedResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	strains := []Strain{}
	for _, hit := range response.Hits.Hits {
		strains = append(strains, hit.Source)
	}
	return strains, nil
}

func (repo FastaRepositoryElasticSearch) MarkAsDownloaded(strainId int) error {
	req := esapi.UpdateRequest{
		Index:      repo.index,
		DocumentID: strconv.Itoa(strainId),
		Body:       strings.NewReader(`{"doc": {"downloaded": true}}`),
	}
	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
