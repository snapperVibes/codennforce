package main

import (
	"encoding/json"
	//"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

/* ------------------------------------------------------------------------------------
Domain (Business Logic)
*/

//
type ParcelService interface {
	Insert(parcelId string) error
	Lookup(parcelId string) io.Writer
}

type wprdcResponse struct {
	// Exported fields?
	Help    string      `json:"help"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Err     interface{} `json:"error"`
}

type recordsType []interface{}
type recordType map[string]interface{}

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

func scrapeRealEstatePortal(parcelID string) []byte {
	repUrl := fmt.Sprintf("http://www2.alleghenycounty.us/RealEstate/Building.aspx?ParcelID=%s", parcelId)
	r, err := http.Get(repUrl)
	if err !=nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = r.Body.Close()
	return bytes
}

func main() {
	var records recordsType
	var record recordType

	municipalityCode := "111"
	wprdcData := downloadDataFromWprdc(municipalityCode)
	//rawWprdcData := string(wprdcData)
	var response wprdcResponse
	jsonErr := json.Unmarshal(wprdcData, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	records = response.Result.(map[string]interface{})["records"].([]interface{})

	for _, _record := range records {
		record = _record.(map[string]interface {})
		parcelId := record["PARID"].(string)
		b := scrapeRealEstatePortal(parcelId)
	}

}

