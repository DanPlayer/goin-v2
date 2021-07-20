package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

type EsOptions struct {
	Addr     string // A list of Elasticsearch nodes to use
	Username string // Username for HTTP Basic Authentication.
	Password string // Password for HTTP Basic Authentication.
}

func NewEs(options EsOptions) *elasticsearch.Client {
	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{options.Addr},
			Username:  options.Username,
			Password:  options.Password,
		})
	if err != nil {
		fmt.Printf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}
	fmt.Println(res)
	return es
}
