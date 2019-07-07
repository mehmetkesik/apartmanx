package apartmanx

import (
	"database/sql"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
	"fmt"
)

type Menu struct {
	Id         int
	Ad         string
	Link       string
	Secenekler []MenuSecenek
}

type MenuSecenek struct {
	Id     int
	Ad     string
	MenuId int
	Link   string
}

func connectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	return db
}

func getMenu() []Menu {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	rows, err := db.Query("select * from menu order by id asc")
	hata(err)
	var menuler []Menu
	for rows.Next() {
		var id int
		var ad, link string
		rows.Scan(&id, &ad, &link)
		var menu Menu
		menu.Id = id
		menu.Ad = ad
		menu.Link = link
		rows2, err := db.Query("select * from menusecenek where menuid = ? order by id asc", menu.Id)
		hata(err)
		var menusecenek []MenuSecenek
		for rows2.Next() {
			var id, menuid int
			var ad, link string
			rows2.Scan(&id, &ad, &menuid, &link)
			var secenek MenuSecenek
			secenek.Id = id
			secenek.Ad = ad
			secenek.MenuId = menuid
			secenek.Link = link
			menusecenek = append(menusecenek, secenek)
		}
		menu.Secenekler = menusecenek
		menuler = append(menuler, menu)
	}
	return menuler
}

func adBul(link string, menu []Menu) string {
	veri := strings.Split(link, "/")
	for _, value := range menu {
		if strings.Trim(veri[1], " ") == value.Link {
			for _, j := range value.Secenekler {
				if j.Link == strings.Trim(veri[2], " ") {
					return j.Ad
				}
			}
		}
	}
	return ""
}

func oturumKontrol(w http.ResponseWriter, r *http.Request) bool {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	var oturum string
	err = db.QueryRow("select value from diger where key = 'oturum'").Scan(&oturum)
	hata(err)
	if oturum == "0" {
		return false
	}
	return true
}

func siteleriGetir() []Site {
	var site []Site
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	rows, err := db.Query("select * from site")
	hata(err)
	for rows.Next() {
		var s Site
		rows.Scan(&s.Id, &s.Sitead, &s.Yoneticiad, &s.Foto, &s.Hakkinda, &s.Telefon, &s.Ceptelefon, &s.Adres, &s.Vergidaire, &s.Vergino, &s.Zaman)
		site = append(site, s)
	}
	return site
}

func aktifSiteGetir() string {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	var aktifsite string
	err = db.QueryRow("select value from diger where key = 'aktifsite'").Scan(&aktifsite)
	hata(err)
	return aktifsite
}

func aktifSiteAd() string {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	id, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var ad string
	err = db.QueryRow("select ad from site where id = ?", id).Scan(&ad)
	hata(err)
	return ad
}

func aktifSiteBilgi() Site {
	var site Site
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	id, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	err = db.QueryRow("select * from site where id = ?", id).Scan(&site.Id, &site.Sitead, &site.Yoneticiad, &site.Foto,
		&site.Hakkinda, &site.Telefon, &site.Ceptelefon, &site.Adres, &site.Vergidaire, &site.Vergino, &site.Zaman)
	hata(err)
	return site
}

func dairelerGetir() []Daire {
	var daireler []Daire
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	id, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	rows, err := db.Query("select * from daire where siteid = ?", id)
	hata(err)
	for rows.Next() {
		var daire Daire
		rows.Scan(&daire.Id, &daire.Siteid, &daire.Blok, &daire.Daireno, &daire.MulksahibiId, &daire.KiraciId,
			&daire.TipiId, &daire.Zaman, &daire.Ekstra)

		var tutsoyad string
		db.QueryRow("select ad,soyad from kisi where id = ?", daire.MulksahibiId).Scan(&daire.Mulksahibi, &tutsoyad)
		daire.Mulksahibi += " " + tutsoyad
		if daire.KiraciId != -1 {
			db.QueryRow("select ad,soyad from kisi where id = ?", daire.KiraciId).Scan(&daire.Kiraci, &tutsoyad)
			daire.Kiraci += " " + tutsoyad
		}

		zamanint, err := strconv.ParseInt(daire.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		daire.Zaman = sdizi[0] + " " + sdizi[1]
		daireler = append(daireler, daire)
	}
	return daireler
}

func daireHareketlerGetir(blok, daireno string) []DaireHareket {
	var dairehareketler []DaireHareket
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	id, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	rows, err := db.Query("select * from dairehareket where siteid = ? and blok = ? and daireno = ? order by tarih desc", id, blok, daireno)
	hata(err)
	for rows.Next() {
		var dairehareket DaireHareket
		var s string
		rows.Scan(&dairehareket.Id, &dairehareket.Siteid, &dairehareket.Blok, &dairehareket.Daireno, &dairehareket.MulksahibiId,
			&dairehareket.KiraciId, &s, &dairehareket.Tarih, &dairehareket.Zaman)

		var tutsoyad string
		db.QueryRow("select ad,soyad from kisi where id = ?", dairehareket.MulksahibiId).Scan(&dairehareket.Mulksahibi, &tutsoyad)
		dairehareket.Mulksahibi += " " + tutsoyad
		if dairehareket.KiraciId != -1 {
			db.QueryRow("select ad,soyad from kisi where id = ?", dairehareket.KiraciId).Scan(&dairehareket.Kiraci, &tutsoyad)
			dairehareket.Kiraci += " " + tutsoyad
		}

		if s != "-1" {
			degisen := strings.Split(s, ",")
			for _, j := range degisen {
				i, err := strconv.Atoi(j)
				hata(err)
				dairehareket.Degisen[i] = true
			}
		}
		zamanint, err := strconv.ParseInt(dairehareket.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		dairehareket.Zaman = sdizi[0] + " " + sdizi[1]

		dairehareketler = append(dairehareketler, dairehareket)
	}
	return dairehareketler
}

func tipiGetir() []Tipi {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	rows, err := db.Query("select * from mulktip")
	hata(err)
	var tipler []Tipi
	for rows.Next() {
		var tip Tipi
		rows.Scan(&tip.Id, &tip.Ad)
		tipler = append(tipler, tip)
	}
	return tipler
}

func kisilerGetir() []Kisi {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var kisiler []Kisi
	rows, err := db.Query("select * from kisi where siteid = ? order by id desc", siteid)
	hata(err)
	for rows.Next() {
		var kisi Kisi
		err = rows.Scan(&kisi.Id, &kisi.Siteid, &kisi.Ad, &kisi.Soyad, &kisi.Mail, &kisi.Telefon, &kisi.Cinsiyet, &kisi.Yas, &kisi.Zaman)
		hata(err)
		zamanint, err := strconv.ParseInt(kisi.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		kisi.Zaman = sdizi[0] + " " + sdizi[1]
		kisiler = append(kisiler, kisi)
	}
	return kisiler
}

func kisiBilgiAra(aramaveri string) []Kisi {
	aramasaf := aramaveri
	aramaveri = "%" + aramaveri + "%"
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var kisiler []Kisi

	var rows *sql.Rows
	if _, err := strconv.Atoi(aramasaf); err == nil {
		rows, err = db.Query("select * from kisi where siteid = ? and id = ? order by id desc",
			siteid, aramasaf)
	} else {
		rows, err = db.Query("select * from kisi where siteid = ? and (ad like ? or soyad like ? or mail like ? or telefon like ?) order by id desc",
			siteid, aramaveri, aramaveri, aramaveri, aramaveri)
	}
	hata(err)
	for rows.Next() {
		var kisi Kisi
		err = rows.Scan(&kisi.Id, &kisi.Siteid, &kisi.Ad, &kisi.Soyad, &kisi.Mail, &kisi.Telefon, &kisi.Cinsiyet, &kisi.Yas, &kisi.Zaman)
		hata(err)
		zamanint, err := strconv.ParseInt(kisi.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		kisi.Zaman = sdizi[0] + " " + sdizi[1]
		kisiler = append(kisiler, kisi)
	}
	return kisiler
}

/////////////////////////////////
func firmalarGetir() []Firma {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var firmalar []Firma
	rows, err := db.Query("select * from firma where siteid = ? order by id desc", siteid)
	hata(err)
	for rows.Next() {
		var firma Firma
		err = rows.Scan(&firma.Id, &firma.Siteid, &firma.Ad, &firma.Sahibad, &firma.Mail, &firma.Telefon, &firma.Adres, &firma.Zaman)
		hata(err)
		zamanint, err := strconv.ParseInt(firma.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		firma.Zaman = sdizi[0] + " " + sdizi[1]
		firmalar = append(firmalar, firma)
	}
	return firmalar
}

func firmaBilgiAra(aramaveri string) []Firma {
	aramasaf := aramaveri
	aramaveri = "%" + aramaveri + "%"
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var firmalar []Firma

	var rows *sql.Rows
	if _, err := strconv.Atoi(aramasaf); err == nil {
		rows, err = db.Query("select * from firma where siteid = ? and id = ? order by id desc",
			siteid, aramasaf)
	} else {
		rows, err = db.Query("select * from firma where siteid = ? and (ad like ? or sahibad like ? or mail like ? or telefon like ?) order by id desc",
			siteid, aramaveri, aramaveri, aramaveri, aramaveri)
	}
	hata(err)
	for rows.Next() {
		var firma Firma
		err = rows.Scan(&firma.Id, &firma.Siteid, &firma.Ad, &firma.Sahibad, &firma.Mail, &firma.Telefon, &firma.Adres, &firma.Zaman)
		hata(err)
		zamanint, err := strconv.ParseInt(firma.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		firma.Zaman = sdizi[0] + " " + sdizi[1]
		firmalar = append(firmalar, firma)
	}
	return firmalar
}

///////////////////////////////////
func gelirgiderlerGetir() []GelirGider {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var gelirgiderler []GelirGider
	rows, err := db.Query("select * from gelirgider where siteid = ? order by id desc", siteid)
	hata(err)
	for rows.Next() {
		var gelirgider GelirGider
		err = rows.Scan(&gelirgider.Id, &gelirgider.Siteid, &gelirgider.Ad, &gelirgider.Zaman)
		hata(err)
		zamanint, err := strconv.ParseInt(gelirgider.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		gelirgider.Zaman = sdizi[0] + " " + sdizi[1]
		gelirgiderler = append(gelirgiderler, gelirgider)
	}
	return gelirgiderler
}

func gelirgiderBilgiAra(aramaveri string) []GelirGider {
	aramasaf := aramaveri
	aramaveri = "%" + aramaveri + "%"
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var gelirgiderler []GelirGider

	var rows *sql.Rows
	if _, err := strconv.Atoi(aramasaf); err == nil {
		rows, err = db.Query("select * from gelirgider where siteid = ? and id = ? order by id desc",
			siteid, aramasaf)
	} else {
		rows, err = db.Query("select * from gelirgider where siteid = ? and (ad like ?) order by id desc",
			siteid, aramaveri)
	}
	hata(err)
	for rows.Next() {
		var gelirgider GelirGider
		err = rows.Scan(&gelirgider.Id, &gelirgider.Siteid, &gelirgider.Ad, &gelirgider.Zaman)
		hata(err)
		zamanint, err := strconv.ParseInt(gelirgider.Zaman, 10, 64)
		hata(err)
		sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
		gelirgider.Zaman = sdizi[0] + " " + sdizi[1]
		gelirgiderler = append(gelirgiderler, gelirgider)
	}
	return gelirgiderler
}

func borcyaz(tutx map[string]string) {
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	zaman := strconv.FormatInt(time.Now().Unix(), 10)
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	_, err = db.Exec("insert into borclandirma (siteid,blok,daireno,miktar,kime,zamangrup,tarih,gelirgiderid,firmaid,'not',zaman) values "+
		"(?,?,?,?,?,?,?,?,?,?,?)", siteid, tutx["blok"], tutx["daireno"], tutx["miktar"], tutx["kime"], tutx["zamangrup"],
		tutx["tarih"+tutx["zamangrup"]], tutx["grup"], tutx["firma"], tutx["not"], zaman)
	hata(err)
}

func borclarGetir(tumborclar bool, blok, dairenostr string) []Borclandirma {
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var sorgustr string
	if tumborclar {
		sorgustr = "select * from borclandirma where siteid = " + strconv.Itoa(siteid) + " order by id desc"
	} else {
		_, err := strconv.Atoi(dairenostr)
		hata(err)
		sorgustr = "select * from borclandirma where siteid = " + strconv.Itoa(siteid) + " and blok = '" + blok + "' and daireno = " + dairenostr +
			" order by id desc"
	}
	var borclandirmalar []Borclandirma
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	rows, err := db.Query(sorgustr)
	if err != nil {
		fmt.Println(err)
	}
	hata(err)

	for rows.Next() {
		var bo Borclandirma
		err = rows.Scan(&bo.Id, &bo.Siteid, &bo.Blok, &bo.Daireno, &bo.Miktar, &bo.Kime, &bo.Zamangrup, &bo.Tarih, &bo.Gelirgiderid,
			&bo.Firmaid, &bo.Not, &bo.Zaman)
		hata(err)

		//Odemeleri getirir.
		rows2, err := db.Query("select * from odeme where borcid = ? and siteid = ?", bo.Id, siteid)
		hata(err)
		for rows2.Next() {
			var o Odeme
			err = rows2.Scan(&o.Id, &o.Siteid, &o.Borcid, &o.Kisiid, &o.Miktar, &o.Tarih, &o.Zaman)
			hata(err)
			var soyad string
			err = db.QueryRow("select ad,soyad from kisi where id = ? and siteid = ?", o.Kisiid, siteid).Scan(&o.Kisiad, &soyad)
			hata(err)
			o.Kisiad += " " + soyad
			para, err := strconv.ParseFloat(o.Miktar, 64)
			hata(err)
			bo.OdemeToplam += para
			bo.Odemeler = append(bo.Odemeler, o)
		}

		err = db.QueryRow("select ad from gelirgider where id = ?", bo.Gelirgiderid).Scan(&bo.Gelirgiderad)
		hata(err)
		if bo.Firmaid != -1 {
			err = db.QueryRow("select ad from firma where id = ?", bo.Firmaid).Scan(&bo.Firmaad)
			hata(err)
		}
		borclandirmalar = append(borclandirmalar, bo)
	}
	return borclandirmalar
}

func odemeyaz(tutx map[string]string) int64 {
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	zaman := strconv.FormatInt(time.Now().Unix(), 10)
	borcid, err := strconv.Atoi(tutx["borcid"])
	hata(err)
	kisiid, err := strconv.Atoi(tutx["kisiid"])
	hata(err)
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	result, err := db.Exec("insert into odeme (siteid,borcid,kisiid,miktar,tarih,zaman) values "+
		"(?,?,?,?,?,?)", siteid, borcid, kisiid, tutx["miktar"], tutx["tarih"], zaman)
	hata(err)
	sonuc, err := result.LastInsertId()
	hata(err)
	return sonuc
}

func dairekontrol(blok, daireno string) bool {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	var id int
	err = db.QueryRow("select id from daire where blok = ? and daireno = ?", blok, daireno).Scan(&id)
	if err != nil {
		return false
	} else {
		return true
	}
}
