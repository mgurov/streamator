package main

import (
	"gopkg.in/sohlich/elogrus.v1"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v3"
)

func newEsHook() (logrus.Hook, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}	
	return elogrus.NewElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
}