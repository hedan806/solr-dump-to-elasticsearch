package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Response is a SOLR response.
type Response struct {
	Header struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
		Params struct {
			Query      string `json:"q"`
			CursorMark string `json:"cursorMark"`
			Sort       string `json:"sort"`
			Rows       string `json:"rows"`
		} `json:"params"`
	} `json:"header"`
	Response struct {
		NumFound int               `json:"numFound"`
		Start    int               `json:"start"`
		Docs     []json.RawMessage `json:"docs"` // dependent on SOLR schema
	} `json:"response"`
	NextCursorMark string `json:"nextCursorMark"`
}

// PrependSchema http, if missing.
func PrependSchema(s string) string {
	if !strings.HasPrefix(s, "http") {
		return fmt.Sprintf("http://%s", s)
	}
	return s
}
func ExportOne(wg *sync.WaitGroup, data chan<- string, server string, sortField string) {
	defer wg.Done()
	query := "*:*"
	rows := 10
	sort := sortField + " asc"
	wt := "json"
	server = PrependSchema(server)

	v := url.Values{}

	v.Set("q", query)
	v.Set("sort", sort)
	v.Set("rows", fmt.Sprintf("%d", rows))

	v.Set("wt", wt)
	v.Set("cursorMark", "*")

	link := fmt.Sprintf("%s/query?%s", server, v.Encode())
	log.Println(link)
	resp, err := http.Get(link)
	if err != nil {
		log.Fatalf("http: %s", err)
	}
	var response Response
	switch wt {
	case "json":
		// invalid character '\r' in string literal
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&response); err != nil {
			log.Fatalf("decode: %s", err)
		}
	default:
		log.Fatalf("wt=%s not implemented", wt)
	}
	// We do not defer, since we hard-exit on errors anyway.
	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}
	for _, doc := range response.Response.Docs {
		data <- string(doc)
	}
}

func Export(wg *sync.WaitGroup, data chan<- string, server string, verbose bool, sortField string) {
	defer wg.Done()
	query := "*:*"
	rows := 100
	sort := sortField + " asc"
	wt := "json"

	flag.Parse()

	server = PrependSchema(server)

	v := url.Values{}

	v.Set("q", query)
	v.Set("sort", sort)
	v.Set("rows", fmt.Sprintf("%d", rows))

	v.Set("wt", wt)
	v.Set("cursorMark", "*")

	var total int

	for {
		link := fmt.Sprintf("%s/query?%s", server, v.Encode())
		if verbose {
			log.Println(link)
		}
		resp, err := http.Get(link)
		if err != nil {
			log.Fatalf("http: %s", err)
		}
		var response Response
		switch wt {
		case "json":
			// invalid character '\r' in string literal
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&response); err != nil {
				log.Fatalf("decode: %s", err)
			}
		default:
			log.Fatalf("wt=%s not implemented", wt)
		}
		// We do not defer, since we hard-exit on errors anyway.
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
		for _, doc := range response.Response.Docs {
			data <- string(doc)
		}
		total += len(response.Response.Docs)
		if verbose {
			log.Printf("fetched %d docs", total)
		}
		if response.NextCursorMark == v.Get("cursorMark") {
			break
		}
		v.Set("cursorMark", response.NextCursorMark)
	}
	if verbose {
		log.Printf("fetched %d docs", total)
	}
}
