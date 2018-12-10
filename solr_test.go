package main

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"testing"
)

func TestExport(t *testing.T) {
	//Export(url, true, "id")
}

func TestExportOne(t *testing.T) {
	done := make(chan bool)
	data := make(chan string)
	var count int32

	go func() {
		for {
			select {
			case doc := <-data:
				docMap := make(map[string]interface{})
				json.Unmarshal([]byte(doc), &docMap)
				log.Println(docMap)
				atomic.AddInt32(&count, 1)
			}
		}
	}()

	<-done

	log.Println("count", count)
}
