package fasta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type FastaRepository interface {
	NotDownloaded(from int) ([]Strain, error)
	MarkAsDownloaded(strainId int) error

	GetByStrainId(strainId int) (*Strain, error)
	FindWithoutIntegronResult(from int) ([]Strain, error)
	AddIntegrons(strainId int, integrons []string) error
}

type DownloadedRequest struct {
	Source struct {
		Downloaded bool `json:"downloaded"`
	} `json:"source"`
}

type Strain struct {
	Id         int `json:"id"`
	AssemblyId int `json:"best_assembly"`
}

type GetByStrainIdResponse struct {
	Source Strain `json:"_source"`
}

type NotDownloadedResponse struct {
	Hits struct {
		Hits []struct {
			Source Strain `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type fastaRepositoryElasticSearch struct {
	index  string
	client *elasticsearch.Client
}

func NewRepository(index string, client *elasticsearch.Client) fastaRepositoryElasticSearch {
	repository := fastaRepositoryElasticSearch{}
	repository.index = index
	repository.client = client
	return repository
}

func (repo fastaRepositoryElasticSearch) NotDownloaded(from int) ([]Strain, error) {
	reader := strings.NewReader(fmt.Sprintf(`{
		"size": 20,
		"sort": [
		  {
			"id": {
			  "order": "desc"
			}
		  }
		],
		"query": {
		  "bool": {
			"filter": [
			  {
				"range": {
				  "id": {
					"lt": %d
				  }
				}
			  }
			],
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
	  }
	  `, from))

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
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	strains := []Strain{}
	for _, hit := range response.Hits.Hits {
		strains = append(strains, hit.Source)
	}
	return strains, nil
}

func (repo fastaRepositoryElasticSearch) FindWithoutIntegronResult(from int) ([]Strain, error) {
	reader := strings.NewReader(fmt.Sprintf(`{
		"size": 20,
		"sort": [
		  {
			"id": {
			  "order": "desc"
			}
		  }
		],
		"query": {
		  "bool": {
			"filter": [
			  {
				"range": {
				  "id": {
					"lt": %d
				  }
				}
			  }
			],
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
				  "integron_finder": {
					"value": true
				  }
				}
			  },
			  {
				"term": {
				  "downloaded": {
					"value": false
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
	  }
	  `, from))

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
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	strains := []Strain{}
	for _, hit := range response.Hits.Hits {
		strains = append(strains, hit.Source)
	}
	return strains, nil
}

func (repo fastaRepositoryElasticSearch) GetByStrainId(strainId int) (*Strain, error) {
	req := esapi.GetRequest{
		Index:      repo.index,
		DocumentID: strconv.Itoa(strainId),
	}
	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("strain not found")
	}

	strainResponse := new(GetByStrainIdResponse)
	err = json.NewDecoder(res.Body).Decode(&strainResponse)
	if err != nil {
		return nil, err
	}

	return &strainResponse.Source, nil
}

func (repo fastaRepositoryElasticSearch) MarkAsDownloaded(strainId int) error {
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

func (repo fastaRepositoryElasticSearch) AddIntegrons(strainId int, integrons []string) error {
	integronsAsArray, err := json.Marshal(integrons)

	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      repo.index,
		DocumentID: strconv.Itoa(strainId),
		Body:       strings.NewReader(fmt.Sprintf(`{"doc": {"integron_finder": true, "integrons": %s}}`, string(integronsAsArray))),
	}
	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
