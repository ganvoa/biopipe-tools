package platform

import "github.com/elastic/go-elasticsearch/v7"

func NewClient(url string, username string, password string) (*elasticsearch.Client, error) {

	config := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
	}

	client, err := elasticsearch.NewClient(config)

	if err != nil {
		return nil, err
	}

	return client, nil
}
