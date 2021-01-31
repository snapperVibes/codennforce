package main

import (
	"encoding/json"
	"io"
	"log"
)

/* ------------------------------------------------------------------------------------
Domain (Business Logic)
*/

//
type ParcelService interface {
	Insert(parcel_id string) error
	Lookup(parcel_id string) io.Writer
}



func main() {
	municipality_code := "111"
	wprdcData := downloadDataFromWprdc(municipality_code)
	//rawWprdcData := string(wprdcData)
	var theData wprdcResponse
	jsonErr := json.Unmarshal(wprdcData, &theData)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
