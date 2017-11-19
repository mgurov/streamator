package main

import (
	"time"
	"fmt"
	"gopkg.in/sohlich/elogrus.v1"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v3"
)

func newEsHook(elasticURL, elasticReportHost, elasticIndex string) (logrus.Hook, error) {
	client, err := elastic.NewClient(elastic.SetURL(elasticURL), elastic.SetSniff(false), elastic.SetHealthcheck(false))
	if err != nil {
		return nil, fmt.Errorf("Creating es client: %s", err)
	}
	var hook *elogrus.ElasticHook
	hook, err = elogrus.NewElasticHook(client, elasticReportHost, logrus.DebugLevel, elasticIndex)
	for hook == nil || err != nil {
		if (err != nil) {
			fmt.Println("Retrying elastic hook creation after:", err)
			time.Sleep(1 * time.Second)
		}
		hook, err = elogrus.NewElasticHook(client, elasticReportHost, logrus.DebugLevel, elasticIndex)
	}
	return hook, err
}