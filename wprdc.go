// Package containing structures and functions dealing with the Western Pennsylvania Regional Data Center
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type wprdcResponse struct {
	// Exported fields?
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  interface{} `json:"result"`
	Err     interface{} `json:"error"`
}


func downloadDataFromWprdc(municipalityCode string) []byte {
	_ = municipalityCode
	// Todo: Ask an expert: Should I use %q for municipality_code (a single quoted character literal)
	unescapedQuery := fmt.Sprintf(`SELECT * FROM "518b583f-7cc8-4f60-94d0-174cc98310dc" WHERE "MUNICODE" = '%s'`, municipalityCode)
	escapedQuery := url.QueryEscape(unescapedQuery)
	wprdcUrl := "https://data.wprdc.org/api/3/action/datastore_search_sql?sql=" + escapedQuery
	response, err := http.Get(wprdcUrl)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = response.Body.Close()
	return bytes
}
