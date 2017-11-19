package main

import (
	"gopkg.in/sohlich/elogrus.v1"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v3"
)

func newEsHook(elasticURL, elasticReportHost, elasticIndex string) (logrus.Hook, error) {
	client, err := elastic.NewClient(elastic.SetURL(elasticURL), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}	
	return elogrus.NewElasticHook(client, elasticReportHost, logrus.DebugLevel, elasticIndex)
}