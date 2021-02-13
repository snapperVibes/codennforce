package s

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

//import (
//	"bytes"
//	"fmt"
//	_ "github.com/lib/pq"
//	"golang.org/x/net/html"
//	"io"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"net/url"
//	"regexp"
//)
//
//// Todo: Get rid of ServerName, ParcelIDText, MuniText, AddressText, OwnerText
////type generalPage struct {
////	Time          [][]byte
////	ServerName    [][]byte
////	ParcelIDText  [][]byte
////	ParcelID      [][]byte
////	MuniText      [][]byte
////	Muni          [][]byte
////	AddressText   [][]byte
////	Address       [][]byte
////	OwnerText     [][]byte
////	Owner         [][]byte
////	School        [][]byte
////	Neighbor      [][]byte
////	Tax           [][]byte
////	OwnerCode     [][]byte
////	State         [][]byte
////	RecDate       [][]byte
////	Use           [][]byte
////	SaleDate      [][]byte
////	Homestead     [][]byte
////	SalePrice     [][]byte
////	Farmstead     [][]byte
////	DeedBook      [][]byte
////	CleanGreen    [][]byte
////	DeedPage      [][]byte
////	Abatement     [][]byte
////	Lot           [][]byte
////	FullLand      [][]byte
////	CountyLand    [][]byte
////	FullBuild     [][]byte
////	CountyBuild   [][]byte
////	FullTot       [][]byte
////	CountyTot     [][]byte
////	FullLand12    [][]byte
////	CountyLand12  [][]byte
////	fullBuild12   [][]byte // Todo: Figure out how to deal with the lowercase
////	CountyBuild12 [][]byte
////	FullTot12     [][]byte
////	CountyTot12   [][]byte
////	ChangeMail    [][]byte
////}
//
type wprdcResponse struct {
	// Exported fields?
	Help    string      `json:"help"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Err     interface{} `json:"error"`
}

//
//func handleErr(err error) {
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
func parseGeneralPage(r io.ReadCloser) map[string][][]byte {
	defer r.Close()
	generalPage := make(map[string][][]byte)
	generalPageLabel := regexp.MustCompile(`_?lbl(?P<label>.+)`)
	z := html.NewTokenizer(r)
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		// Check if token represents labeled data
		id, val, _ := z.TagAttr()
		compared := bytes.Compare(id, []byte("id"))
		if compared != 0 {
			continue
		}
		m := generalPageLabel.FindSubmatch(val)
		if m == nil {
			continue
		}
		// Todo: Create a script to programmatically create these structs for me
		lbl := string(m[1])

		//println(string(lbl))
		for {

			tt = z.Next()
			if tt == html.TextToken {
				var buf bytes.Buffer
				buf.Write(z.Raw())
				v := generalPage[lbl]
				v = append(v, buf.Bytes())
				generalPage[lbl] = v
				//println(string(v[0]))
				continue
			}
			if tt == html.StartTagToken {
				continue
			}
			break
		}
	}
	return generalPage
}

//
func downloadDataFromWprdc(municipalityCode string) []byte {
	// Todo: Ask an expert: Should I use %q for municipality_code (a single quoted character literal)
	unescapedQuery := fmt.Sprintf(`SELECT * FROM "518b583f-7cc8-4f60-94d0-174cc98310dc" WHERE "MUNICODE" = '%s'`, municipalityCode)
	escapedQuery := url.QueryEscape(unescapedQuery)
	wprdcUrl := "https://data.wprdc.org/api/3/action/datastore_search_sql?sql=" + escapedQuery
	response, err := http.Get(wprdcUrl)
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = response.Body.Close()
	return b
}

type recordsType []interface{}
type recordType map[string]interface{}

//func main() {
//
//	var records recordsType
//	//urls := make(chan recordType, 6)
//
//	municipalityCode := "111"
//	wprdcData := downloadDataFromWprdc(municipalityCode)
//	//rawWprdcData := string(wprdcData)
//	var response wprdcResponse
//	jsonErr := json.Unmarshal(wprdcData, &response)
//	if jsonErr != nil {
//		log.Fatal(jsonErr)
//	}
//	records = response.Result.(map[string]interface{})["records"].([]interface{})
//
//	for i, _record := range records {
//		record := _record.(map[string]interface{})
//		parcelId := record["PARID"].(string)
//
//		r, err := http.Get(fmt.Sprintf("http://www2.alleghenycounty.us/RealEstate/GeneralInfo.aspx?ParcelID=%s", parcelId))
//		handleErr(err)
//		myMap := parseGeneralPage(r.Body)
//		fmt.Printf("%d\t%s\n", i, myMap)
//	}
//}

func connectToDB() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "cogdb"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}

func insertScraped(db *sql.DB, m map[string][][]byte, url string) error {
	//parcelid	ParcelID
	//propertyaddress	Address
	//municipality	Muni
	//ownername	Owner
	//ownermailing	ChangeMail
	insertSQL := `
INSERT INTO realestateportal (url, address, municipality, owner, changemail)
VALUES ($1, $2, $3, $4, $5);`
	_, err := db.Exec(insertSQL, url, pq.Array(m["Address"]), pq.Array(m["Muni"]), pq.Array(m["Owner"]), pq.Array(m["ChangeMail"]))
	return err
}

func main() {

	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// Select municodes
	rows, err := db.Query("SELECT municode FROM municipality WHERE municode != 999;")
	if err != nil {
		log.Fatal(err)
	}
	var municodes []int
	for rows.Next() {
		var municode int
		err := rows.Scan(&municode)
		if err != nil {
			log.Fatal(err)
		}
		municodes = append(municodes, municode)
	}
	if municodes == nil {
		log.Fatal(municodes)
	}
	fmt.Println(municodes)

	// Download Data from WPRDC
	b := downloadDataFromWprdc(strconv.Itoa(municodes[0]))
	var response wprdcResponse
	jsonErr := json.Unmarshal(b, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	var records recordsType
	records = response.Result.(map[string]interface{})["records"].([]interface{})

	for i, _record := range records {
		record := _record.(map[string]interface{})
		parcelId := record["PARID"].(string)
		portalUrl := fmt.Sprintf("http://www2.alleghenycounty.us/RealEstate/GeneralInfo.aspx?ParcelID=%s", parcelId)
		r, err := http.Get(portalUrl)
		if err != nil {
			log.Fatal(err)
		}
		myMap := parseGeneralPage(r.Body)
		//fmt.Printf("%d\t%s\n", i, myMap)

		err = insertScraped(db, myMap, portalUrl)
		if err != nil {
			log.Fatal(err)
		}
		println(i)
	}
}
