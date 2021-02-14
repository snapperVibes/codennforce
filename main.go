package main

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
	"strings"
)

// One or more space or non-breaking space

var spaces = regexp.MustCompile(`[\x{0020}\x{00A0}]+`)

type scrapedHTML [][]byte
type record = []interface{}
type wprdcResponse struct {
	Help    string      `json:"help"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Err     interface{} `json:"error"`
}

// Todo: Add confirmation
const SnapScrapeId = 20

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	municodes, err := selectMunicodes(db)
	if err != nil {
		log.Fatal(err)
	}

	// Download Data from WPRDC
	b := downloadDataFromWprdc(strconv.Itoa(municodes[0]))
	var response wprdcResponse
	jsonErr := json.Unmarshal(b, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	var records record
	records = response.Result.(map[string]interface{})["records"].(record)

	for i, _record := range records {
		m := _record.(map[string]interface{})
		parcelId := m["PARID"].(string)
		portalUrl := fmt.Sprintf("http://www2.alleghenycounty.us/RealEstate/GeneralInfo.aspx?ParcelID=%s", parcelId)
		r, err := http.Get(portalUrl)
		if err != nil {
			log.Fatal(err)
		}
		myMap := parseGeneralPage(r.Body)

		// DATABASE INSERTS
		err = insertRealEstatePortal(db, myMap, portalUrl)
		if err != nil {
			log.Fatal(err)
		}
		parcelKey, err := ensureParcel(db, parcelId)
		if err != nil {
			log.Fatal(err)
		}
		dirtyOwner := myMap["Owner"]
		owner := cleanScrapedHTML(dirtyOwner)
		err = insertOwners(db, parcelKey, owner)
		if err != nil {
			log.Fatal(err)
		}
		dirtyAddress := myMap["Address"]
		address := cleanScrapedHTML(dirtyAddress)
		err = insertAddress(db, parcelKey, address)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d\t%s\t%s\n", i, owner, address)
	}
}

func selectMunicodes(db *sql.DB) ([]int, error) {
	var municodes []int
	rows, err := db.Query("SELECT municode FROM municipality WHERE municode != 999;")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var municode int
		err = rows.Scan(&municode)
		municodes = append(municodes, municode)
	}
	fmt.Println("Municodes selected")
	return municodes, err
}

func parseGeneralPage(r io.ReadCloser) map[string]scrapedHTML {
	defer r.Close()
	generalPage := make(map[string]scrapedHTML)
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

func ensureParcel(db *sql.DB, parcelID string) (int, error) {

	selectSql := "SELECT id FROM parcel WHERE parcelid = $1;"
	row := db.QueryRow(selectSql, parcelID)
	if row.Err() != nil {
		log.Fatal(row.Err())
	}
	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		insertSql := "INSERT INTO parcel (parcelid) VALUES ($1) RETURNING id;"
		row = db.QueryRow(insertSql, parcelID)
		err = row.Scan(&id)
	}
	return id, err
}

func cleanScrapedHTML(scrapedHTML scrapedHTML) (ret []string) {
	for _, _html := range scrapedHTML {
		s := string(_html)
		s = html.UnescapeString(s)
		s = spaces.ReplaceAllString(s, " ")
		s = strings.TrimSpace(s)
		ret = append(ret, s)
	}
	return ret
}

func insertRealEstatePortal(db *sql.DB, m map[string]scrapedHTML, url string) error {
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

func insertOwners(db *sql.DB, key int, cleanedNames []string) error {
	insertSql := "INSERT INTO owner (parcel_id, fullname, bobsource_sourceid) VALUES ($1, $2, $3);"
	_, err := db.Exec(insertSql, key, pq.Array(cleanedNames), SnapScrapeId)
	return err
}

func insertAddress(db *sql.DB, key int, address []string) error {
	insertSql := `
INSERT INTO address (parcel_id, fulladdress, bobsource_sourceid)
VALUES ($1, $2, $3)`
	_, err := db.Exec(insertSql, key, pq.Array(address), SnapScrapeId)
	return err
}
