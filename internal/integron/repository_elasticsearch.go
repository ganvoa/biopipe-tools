package integron

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type integronRepository struct {
	index  string
	client *elasticsearch.Client
}

type IntegronUpdateRequest struct {
	Doc struct {
		IntegronResult  []Integron
		Integron_finder bool
	}
}

func NewRepository(index string, client *elasticsearch.Client) integronRepository {
	repository := integronRepository{}
	repository.index = index
	repository.client = client
	return repository
}

func (repo integronRepository) AddIntegron(strainId int, integrons []Integron) error {

	body := IntegronUpdateRequest{
		Doc: struct {
			IntegronResult  []Integron
			Integron_finder bool
		}{IntegronResult: integrons, Integron_finder: true},
	}

	updateString, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      repo.index,
		DocumentID: strconv.Itoa(strainId),
		Body:       bytes.NewBuffer(updateString),
	}
	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(strconv.Itoa(res.StatusCode))
	}

	return nil
}
