package externaldata

type scrapedHTML [][]byte

type RealEstatePortal struct {
	Farmstead     scrapedHTML `id:"lblFarmstead"`
	Tax           scrapedHTML `id:"lblTax"`
	Abatement     scrapedHTML `id:"lblAbatement"`
	ParcelID      scrapedHTML `id:"BasicInfo1_lblParcelID"`
	FullLand      scrapedHTML `id:"lblFullLand"`
	FullBuild     scrapedHTML `id:"lblFullBuild"`
	CountyTot12   scrapedHTML `id:"lblCountyTot12"`
	SaleDate      scrapedHTML `id:"lblSaleDate"`
	CleanGreen    scrapedHTML `id:"lblCleanGreen"`
	FullTot       scrapedHTML `id:"lblFullTot"`
	State         scrapedHTML `id:"lblState"`
	Neighbor      scrapedHTML `id:"lblNeighbor"`
	ServerName    scrapedHTML `id:"Header1_lblServerName"`
	CountyLand    scrapedHTML `id:"lblCountyLand"`
	fullBuild12   scrapedHTML `id:"lblfullBuild12"`
	Address       scrapedHTML `id:"BasicInfo1_lblAddress"`
	Lot           scrapedHTML `id:"lblLot"`
	CountyLand12  scrapedHTML `id:"lblCountyLand12"`
	OwnerCode     scrapedHTML `id:"lblOwnerCode"`
	Muni          scrapedHTML `id:"BasicInfo1_lblMuni"`
	CountyBuild   scrapedHTML `id:"lblCountyBuild"`
	Use           scrapedHTML `id:"lblUse"`
	DeedBook      scrapedHTML `id:"lblDeedBook"`
	School        scrapedHTML `id:"lblSchool"`
	RecDate       scrapedHTML `id:"lblRecDate"`
	FullLand12    scrapedHTML `id:"lblFullLand12"`
	ChangeMail    scrapedHTML `id:"lblChangeMail"`
	CountyBuild12 scrapedHTML `id:"lblCountyBuild12"`
	Owner         scrapedHTML `id:"BasicInfo1_lblOwner"`
	Homestead     scrapedHTML `id:"lblHomestead"`
	SalePrice     scrapedHTML `id:"lblSalePrice"`
	Time          scrapedHTML `id:"Header1_lblTime"`
	CountyTot     scrapedHTML `id:"lblCountyTot"`
	DeedPage      scrapedHTML `id:"lblDeedPage"`
	FullTot12     scrapedHTML `id:"lblFullTot12"`
}
