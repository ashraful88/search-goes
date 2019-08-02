package search

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/tidwall/gjson"
)

func buildQuery(key, val string) io.Reader {
	var b strings.Builder

	switch key {
	case "q":
		b.WriteString(fmt.Sprintf(searchByKeyword, val))
	case "category":
	case "region":

	default:
		b.WriteString(searchAll)
	}

	return strings.NewReader(b.String())
}

// QuerySearch perform a query search in ES
func QuerySearch(q, category, region string) string {
	se := GetSearchClient()
	// Perform the search request.
	log.Println(category, "category")
	log.Println(region, "region")
	log.Println(buildQuery("q", "myvi"))
	res, err := se.Client.Search(
		se.Client.Search.WithContext(context.Background()),
		se.Client.Search.WithIndex(se.AdsIndexName),
		//se.Client.Search.WithSize(40),

		se.Client.Search.WithBody(buildQuery("q", "myvi")),
		se.Client.Search.WithTrackTotalHits(true),
		se.Client.Search.WithPretty(),
	)

	//log.Println(res)
	if res.IsError() {
	}

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	json := read(res.Body)
	value := gjson.Get(json, "name.last")
	println(value.String())
	return json
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}

const searchAll = `{
	"query" : { "match_all" : {} },
	"size" : 40,
	"sort" : { "list_date" : "desc"}
	}`

const searchByCategoryRegion = `{
	"query": {
	  "bool": {
		"must": [
		  {"match": {
			"category": %q
		  }},
		  {"match": {
			"region": "%q
		  }}
		]
	  }
	}
  }`

const searchByKeyword = `{
	"query" : {
		"multi_match" : {
			"query" : %q,
			"fields" : ["subject^3", "body"],
			"operator" : "and"
		}
	},
	"size" : 40,
	"sort" : [ { "_score" : "desc" }, { "list_date" : "desc" } ]
	}`