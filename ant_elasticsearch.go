package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

var (
	esClient *elastic.Client
	ctx      = context.Background()
)

func Save(url string, data []interface{}) {
	esClient = Client(url)
	bulkService := *esClient.Bulk()
	for _, doc := range data {
		req := elastic.NewBulkIndexRequest().Index("index").Type("doc").Doc(doc).UseEasyJSON(true)
		bulkService.Add(req)
	}
	bulkResponse, err := bulkService.Do(ctx)
	if err != nil {
		log.Println(err)
	}
	indexed := bulkResponse.Items
	log.Println("向es导入了", len(indexed), "条数据")
}

func SaveOne(url string, data interface{}) {
	esClient = Client(url)
	response, err := esClient.Index().Index("index").Type("doc").BodyJson(data).Do(ctx)
	if err != nil {
	}
	indexed := response.Status
	fmt.Println("向es导入了", indexed, "条数据")
}

func Client(url string) *elastic.Client {
	urls := []string{url}
	client, err := elastic.NewClient(elastic.SetURL(urls...))
	if err != nil {
		panic(err)
	}
	return client
}
