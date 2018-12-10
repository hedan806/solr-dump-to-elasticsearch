package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"sync/atomic"
	"testing"
)

func TestSave(t *testing.T) {

	done := make(chan bool)
	data := make(chan string, 10)
	var afters int64
	var failures int64
	var count int32
	client := Client("http://10.20.1.64:9200")
	afterFn := func(executionId int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
		atomic.AddInt64(&afters, 1)
		if err != nil {
			log.Println(err)
			atomic.AddInt64(&failures, 1)
		}
	}
	bulkProcessor, err := client.BulkProcessor().After(afterFn).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			doc, more := <-data
			if more {
				docMap := make(map[string]interface{})
				json.Unmarshal([]byte(doc), &docMap)
				atomic.AddInt32(&count, 1)
				request := elastic.NewBulkIndexRequest().Index("index").Type("doc").Doc(docMap)
				bulkProcessor.Add(request)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()
	close(data)
	<-done
	log.Println("count", count)
	bulkProcessor.Close()
}
