package apartmanx

import (
	"net/http"
	"strings"
	"html/template"
	"strconv"
	"database/sql"
	"time"
)

type Kisi struct {
	Id       int
	Siteid   int
	Ad       string
	Soyad    string
	Mail     string
	Telefon  string
	Cinsiyet string
	Yas      string
	Zaman    string
}

type Firma struct {
	Id      int
	Siteid  int
	Ad      string
	Sahibad string
	Mail    string
	Telefon string
	Adres   string
	Zaman   string
}

type GelirGider struct {
	Id     int
	Siteid int
	Ad     string
	Zaman  string
}

///Kişi Tanımı
func kisiTanimi(w http.ResponseWriter, r *http.Request) {
	type KisiTanimi struct {
		Apartman
		Link           string
		KisiKayitDurum bool
		KisiBilgi      Kisi
		KisiAra        bool
		KisilerBilgi   []Kisi
	}
	var veri KisiTanimi
	link := strings.Split(r.URL.String(), "?")[0]

	veri.Link = link

	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	if r.Method == "POST" {
		if r.FormValue("kisikayit") != "" {
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			_, err = db.Exec("insert into kisi (siteid,ad,soyad,mail,telefon,cinsiyet,yas,zaman) values (?,?,?,?,?,?,?,?)",
				siteid, r.FormValue("ad"), r.FormValue("soyad"), r.FormValue("mail"), r.FormValue("telefon"),
				r.FormValue("cinsiyet"), r.FormValue("yas"), zaman)
			hata(err)
			http.Redirect(w, r, link+"?durum=1", 302)
		} else {
			http.Redirect(w, r, link, 302)
		}
	}

	veri.KisiKayitDurum = false
	if r.URL.Query().Get("durum") == "1" {
		veri.KisiKayitDurum = true
	}

	veri.KisiAra = false
	if r.URL.Query().Get("kisiara") != "" {
		aramaveri := r.URL.Query().Get("kisiara")
		if aramaveri == "tumkisiler" {
			veri.KisilerBilgi = kisilerGetir()
		} else {
			veri.KisilerBilgi = kisiBilgiAra(aramaveri)
		}
		veri.KisiAra = true
	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}

///Firma Tanımı
func firmaTanimi(w http.ResponseWriter, r *http.Request) {
	type FirmaTanimi struct {
		Apartman
		Link            string
		FirmaKayitDurum bool
		FirmaBilgi      Firma
		FirmaAra        bool
		FirmalarBilgi   []Firma
	}
	var veri FirmaTanimi
	link := strings.Split(r.URL.String(), "?")[0]

	veri.Link = link

	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	if r.Method == "POST" {
		if r.FormValue("firmakayit") != "" {
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			_, err = db.Exec("insert into firma (siteid,ad,sahibad,mail,telefon,adres,zaman) values (?,?,?,?,?,?,?)",
				siteid, r.FormValue("ad"), r.FormValue("sahibad"), r.FormValue("mail"), r.FormValue("telefon"),
				r.FormValue("adres"), zaman)
			hata(err)
			http.Redirect(w, r, link+"?durum=1", 302)
		} else {
			http.Redirect(w, r, link, 302)
		}
	}

	veri.FirmaKayitDurum = false
	if r.URL.Query().Get("durum") == "1" {
		veri.FirmaKayitDurum = true
	}

	veri.FirmaAra = false
	if r.URL.Query().Get("firmaara") != "" {
		aramaveri := r.URL.Query().Get("firmaara")
		if aramaveri == "tumfirmalar" {
			veri.FirmalarBilgi = firmalarGetir()
		} else {
			veri.FirmalarBilgi = firmaBilgiAra(aramaveri)
		}
		veri.FirmaAra = true
	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}

///Gelir-Gider Tanımı
func gelirGiderTanimi(w http.ResponseWriter, r *http.Request) {
	type GelirGiderTanimi struct {
		Apartman
		Link                 string
		GelirgiderKayitDurum bool
		GelirgiderBilgi      GelirGider
		GelirgiderAra        bool
		GelirgiderlerBilgi   []GelirGider
	}
	var veri GelirGiderTanimi
	link := strings.Split(r.URL.String(), "?")[0]

	veri.Link = link

	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	if r.Method == "POST" {
		if r.FormValue("gelirgiderkayit") != "" {
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			_, err = db.Exec("insert into gelirgider (siteid,ad,zaman) values (?,?,?)",
				siteid, r.FormValue("ad"), zaman)
			hata(err)
			http.Redirect(w, r, link+"?durum=1", 302)
		} else {
			http.Redirect(w, r, link, 302)
		}
	}

	veri.GelirgiderKayitDurum = false
	if r.URL.Query().Get("durum") == "1" {
		veri.GelirgiderKayitDurum = true
	}

	veri.GelirgiderAra = false
	if r.URL.Query().Get("gelirgiderara") != "" {
		aramaveri := r.URL.Query().Get("gelirgiderara")
		if aramaveri == "tumgelirgiderler" {
			veri.GelirgiderlerBilgi = gelirgiderlerGetir()
		} else {
			veri.GelirgiderlerBilgi = gelirgiderBilgiAra(aramaveri)
		}
		veri.GelirgiderAra = true
	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}
