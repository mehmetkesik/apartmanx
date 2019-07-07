package apartmanx

import (
	"net/http"
	"strings"
	"html/template"
	"strconv"
	"database/sql"
)

type TanimliRapor struct {
	Id     int
	SiteId int
	Ad     string
	Link   string
	Zaman  string
}

func tanimliRaporlar(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type TanımlıRaporlar struct {
		Apartman
		Link            string
		LinkVeri        LinkVeri
		Rapor           Rapor
		AramaDurum      bool
		TanimliRaporlar []TanimliRapor
		SorguId         int
	}

	var veri TanımlıRaporlar
	link := strings.Split(r.URL.String(), "?")[0]
	veri.Link = link
	veri.TanimliRaporlar = tanimliRaporlarGetir()

	if r.URL.Query().Get("x") != "" {
		var err error
		x, err := strconv.Atoi(r.URL.Query().Get("x"))
		hata(err)
		veri.SorguId = x
		veri.AramaDurum = true
		if r.URL.Query().Get("kisidagit") != "" {
			veri.LinkVeri.KisiDagit = true
		}
		veri.LinkVeri.KisiId, err = strconv.Atoi(r.URL.Query().Get("kisiid"))
		if err != nil {
			veri.LinkVeri.KisiId = -1
		}
		veri.LinkVeri.Kime, err = strconv.Atoi(r.URL.Query().Get("kime"))
		if err != nil {
			veri.LinkVeri.Kime = -1
		}
		veri.LinkVeri.Blok = r.URL.Query().Get("blok")
		veri.LinkVeri.DaireNo, err = strconv.Atoi(r.URL.Query().Get("daireno"))
		if err != nil {
			veri.LinkVeri.DaireNo = -1
		}
		veri.LinkVeri.FirmaId, err = strconv.Atoi(r.URL.Query().Get("firmaid"))
		if err != nil {
			veri.LinkVeri.FirmaId = -1
		}
		veri.LinkVeri.GelirGiderId, err = strconv.Atoi(r.URL.Query().Get("gelirgiderid"))
		if err != nil {
			veri.LinkVeri.GelirGiderId = -1
		}
		veri.LinkVeri.BasTarih = r.URL.Query().Get("bastarih")
		veri.LinkVeri.BitTarih = r.URL.Query().Get("bittarih")

		veri.Rapor = ozelRapor(veri.LinkVeri)
	}

	html, err := template.ParseFiles("root"+link+".html", "root/raporlar/rapor/ozelRaporDaire.html",
		"root/raporlar/rapor/ozelRaporKisi.html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}

func tanimliRaporlarGetir() []TanimliRapor {
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	rows, err := db.Query("select * from tanimlirapor where siteid = ?", siteid)
	hata(err)
	var tanimliraporlar []TanimliRapor
	for rows.Next() {
		var tanimlirapor TanimliRapor
		rows.Scan(&tanimlirapor.Id, &tanimlirapor.SiteId, &tanimlirapor.Ad, &tanimlirapor.Link, &tanimlirapor.Zaman)
		tanimliraporlar = append(tanimliraporlar, tanimlirapor)
	}
	return tanimliraporlar
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////// ÖZEL RAPOR HAZIRLA ////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type LinkVeri struct {
	KisiDagit    bool
	KisiId       int
	Blok         string
	DaireNo      int
	Kime         int
	FirmaId      int
	GelirGiderId int
	BasTarih     string
	BitTarih     string
}

func ozelRaporHazirla(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type OzelRaporHazirla struct {
		Apartman
		Link          string
		DaireBlok     []string
		KisilerBilgi  []Kisi
		Firmalar      []Firma
		GelirGiderler []GelirGider
		LinkVeri      LinkVeri
		Rapor         Rapor
		AramaDurum    bool
	}
	var veri OzelRaporHazirla
	link := strings.Split(r.URL.String(), "?")[0]
	veri.Link = link
	veri.DaireBlok = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U", "V", "Y", "Z"}
	veri.KisilerBilgi = kisilerGetir()
	veri.Firmalar = firmalarGetir()
	veri.GelirGiderler = gelirgiderlerGetir()

	if r.URL.Query().Get("rapor") == "ozelrapor" {
		veri.AramaDurum = true
		var err error
		if r.URL.Query().Get("kisidagit") != "" {
			veri.LinkVeri.KisiDagit = true
		}
		veri.LinkVeri.KisiId, err = strconv.Atoi(r.URL.Query().Get("kisiid"))
		if err != nil {
			veri.LinkVeri.KisiId = -1
		}
		veri.LinkVeri.Kime, err = strconv.Atoi(r.URL.Query().Get("kime"))
		if err != nil {
			veri.LinkVeri.Kime = -1
		}
		veri.LinkVeri.Blok = r.URL.Query().Get("blok")
		veri.LinkVeri.DaireNo, err = strconv.Atoi(r.URL.Query().Get("daireno"))
		if err != nil {
			veri.LinkVeri.DaireNo = -1
		}
		veri.LinkVeri.FirmaId, err = strconv.Atoi(r.URL.Query().Get("firmaid"))
		if err != nil {
			veri.LinkVeri.FirmaId = -1
		}
		veri.LinkVeri.GelirGiderId, err = strconv.Atoi(r.URL.Query().Get("gelirgiderid"))
		if err != nil {
			veri.LinkVeri.GelirGiderId = -1
		}
		veri.LinkVeri.BasTarih = r.URL.Query().Get("bastarih")
		veri.LinkVeri.BitTarih = r.URL.Query().Get("bittarih")

		veri.Rapor = ozelRapor(veri.LinkVeri)
	}

	html, err := template.ParseFiles("root"+link+".html", "root/raporlar/rapor/ozelRaporDaire.html",
		"root/raporlar/rapor/ozelRaporKisi.html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}
