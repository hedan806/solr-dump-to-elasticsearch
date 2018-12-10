package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

var (
	wg  sync.WaitGroup
	wgp sync.WaitGroup
)

func Export2ES(consumerCount int, src string, dest string, index string, indexType string, id string) {
	data := make(chan string, 1000)
	var afters int64
	var failures int64
	var count int32
	client := Client(dest)
	afterFn := func(executionId int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
		atomic.AddInt64(&afters, 1)
		if err != nil {
			log.Println(err)
			atomic.AddInt64(&failures, 1)
		}
	}
	bulkProcessor, err := client.BulkProcessor().After(afterFn).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(consumerCount)
	log.Println("consumer:", consumerCount)
	for i := 0; i < consumerCount; i++ {
		go consumer(&wg, data, &count, i, id, index, indexType, bulkProcessor)
	}
	wgp.Add(1)
	go Export(&wgp, data, src, true, id)

	wgp.Wait()
	close(data)
	wg.Wait()
	log.Println("count", count)
	bulkProcessor.Close()
}

func consumer(wg *sync.WaitGroup, datas <-chan string, count *int32, index int, id string, indexName string, indexType string, bulkProcessor *elastic.BulkProcessor) {
	defer wg.Done()
	for doc := range datas {
		docMap := make(map[string]interface{})
		json.Unmarshal([]byte(doc), &docMap)
		atomic.AddInt32(count, 1)
		//log.Println(index, *count, docMap)
		request := elastic.NewBulkIndexRequest().Id(docMap[id].(string)).Index(indexName).Type(indexType).Doc(docMap)
		bulkProcessor.Add(request)
	}
	fmt.Println(index, "done!")
}
