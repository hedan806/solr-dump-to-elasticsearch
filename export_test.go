package main

import "testing"

func TestExport2ES(t *testing.T) {
	Export2ES(100, "http://172.19.105.13:8983/solr/index",
		"http://10.20.1.64:9200", "eefung", "doc", "id")
}
