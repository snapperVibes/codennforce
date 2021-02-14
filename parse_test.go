package main

import (
	"github.com/facebookgo/ensure"
	"testing"
)

// Test cleanString

func TestSingle(t *testing.T) {
	ensure.DeepEqual(t, cleanString("VIBES SNAPPER     "), "VIBES SNAPPER")
}

func TestSingleSpaceAmpersand(t *testing.T) {
	ensure.DeepEqual(t, cleanString("VIBES SNAPPER JR &  MARY R     "), "VIBES SNAPPER JR & MARY R")
}

func TestLong(t *testing.T) {
	ensure.DeepEqual(
		t,
		cleanString("EAST ALLEGHENY SCHOOL DISTRICT EAST   MCKEESPORT BOROUGH &  ALLEGHENY COUNTY   "),
		"EAST ALLEGHENY SCHOOL DISTRICT EAST MCKEESPORT BOROUGH & ALLEGHENY COUNTY")
}

// Todo: Automatically find and fix
func TestBrokenExample(t *testing.T) {
	ensure.DeepEqual(
		t,
		cleanString("CHURCH OF THE BRETHREN OF EAST MC KEESPO   RT   "),
		"CHURCH OF THE BRETHREN OF EAST MC KEESPO RT")
}

//
//func TestDoubleSpaceAmpersand(t *testing.T) {
//	ensure.DeepEqual(
//		t,
//		ParseOwners("LOVING EDGAR  &  ROBNET LOVING     "),
//		[]string{"LOVING EDGAR", "ROBNET LOVING"})
//}
//
//// Words seperated by triple spaces should be kept together
//func TestTripleSpace(t *testing.T) {
//	db, err := connectToDB()
//	if err != nil {log.Fatal(err)}
//	rows, err := db.Query("SELECT owner FROM realestateportal;")
//	if err != nil {log.Fatal(err)}
//	fmt.Printf("%v", rows)
//}

/*
Real places with an exclamation point at the end
LOVING EDGAR  &  ROBNET LOVING     !
BAUR DAVID J &  JUDITH R REVOCABLE LIVING    TRUST   !
DUMBLOSKY MARK A &  BARBARA A REVOCABLE   LIVING TRUST (THE)!
MELVIN J WOJTOWICZ &  BONNIE E ALEXANDER   REVOCABLE LIVING TRUST!
VITEK STEPHANIE MARIE &  JOSEPH M   (H)   !
HOWARD ROBERT PHILLIP III &  MARIE ELAINE    (W)   !
CHURCH OF THE BRETHREN OF EAST MC KEESPO   RT   !
BRADLEY EDWIN S &  DOROTHY #1     !
EAST ALLEGHENY SCHOOL DISTRICT EAST   MCKEESPORT BOROUGH &  ALLEGHENY COUNTY   !
JOHN B PAGANI REVOCABLE AGREEMENT   OF TRUST   !
DUMBLOSKY MARK  A &  BARBARA A REVOCABLE   LIVING TRUST ( THE )   !

        "MELVIN J WOJTOWICZ &  BONNIE E ALEXANDER   REVOCABLE LIVING TRUST   ",
        "DUMBLOSKY MARK  A &  BARBARA A REVOCABLE   LIVING TRUST ( THE )   "

BAUR DAVID J &  JUDITH R REVOCABLE LIVING    TRUST
DUMBLOSKY MARK A &  BARBARA A REVOCABLE   LIVING TRUST (THE)
DUMBLOSKY MARK  A &  BARBARA A REVOCABLE   LIVING TRUST ( THE )
DUMBLOSKY MARK A &  BARBARA A REVOCABLE   LIVING TRUST (THE)
*/

//func TestExample(t *testing.T) {
//	ensure.DeepEqual(
//		t,
//		ParseOwners("BAUR DAVID J &  JUDITH R REVOCABLE LIVING    TRUST   "),
//		[]string{"BAUR DAVID J", "JUDITH R REVOCABLE LIVING", "TRUST"})
//}

// The data is split up automatically: For example,
// http://www2.alleghenycounty.us/RealEstate/GeneralInfo.aspx?ParcelID=0546H00185000000
// COSMOPOLITAN SAVINGS &  LOAN ASSN OF PGH
// After any ampersand is a double space even if it is a single entity
//
//func TestParseMistakenAmpresand(t *testing.T) {
//	ensure.DeepEqual()
//}
