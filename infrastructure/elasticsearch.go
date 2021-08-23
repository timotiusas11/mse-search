package infrastructure

import elastic "github.com/elastic/go-elasticsearch/v8"

func GetESClient() (*elastic.Client, error) {
	client, err := elastic.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return client, nil
}
