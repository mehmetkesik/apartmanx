package apartmanx

import (
	"net/http"
	"html/template"
	"os"
	"strconv"
	"io/ioutil"
	"database/sql"
	"strings"
	"time"
	"fmt"
)

type Site struct {
	Id         int
	Sitead     string
	Yoneticiad string
	Foto       string
	Hakkinda   string
	Telefon    string
	Ceptelefon string
	Adres      string
	Vergidaire string
	Vergino    string
	Zaman      string
}

func siteEkle(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type SiteEkle struct {
		Apartman
		YoneticiAd string
	}
	var veri SiteEkle
	var profil ProfilX
	profil.profilGetir()

	if r.Method == "POST" {
		if r.FormValue("siteekle") != "" {
			var fotoname string
			sitead := r.FormValue("sitead")
			yoneticiad := r.FormValue("yoneticiad")
			hakkinda := r.FormValue("hakkinda")
			telefon := r.FormValue("telefon")
			ceptelefon := r.FormValue("ceptelefon")
			adres := r.FormValue("adres")
			vergidaire := r.FormValue("vergidaire")
			vergino := r.FormValue("vergino")
			_, fotoheader, err := r.FormFile("foto")
			if err == nil {
				fotoname = fotoheader.Filename
				uzanti := fotoname[strings.LastIndex(fotoname, ".")+1:]
				var sira = 1
				for {
					if _, err := os.Stat("root/img/" + strconv.Itoa(sira) + "." + uzanti); err != nil {
						dosya, err := fotoheader.Open()
						hata(err)
						b, err := ioutil.ReadAll(dosya)
						dosya.Close()
						hata(err)
						dosya2, err := os.Create("root/img/" + strconv.Itoa(sira) + "." + uzanti)
						hata(err)
						dosya2.Write(b)
						dosya2.Close()
						fotoname = strconv.Itoa(sira) + "." + uzanti
						break
					}
					sira += 1
				}
			}
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			db, err := sql.Open("sqlite3", "apartmanx.db")
			_, err = db.Exec("insert into site ('ad','yoneticiad','foto','hakkinda','telefon','ceptelefon','adres','vergidairesi','vergino','zaman')"+
				" values (?,?,?,?,?,?,?,?,?,?)", sitead, yoneticiad, fotoname, hakkinda, telefon, ceptelefon, adres, vergidaire, vergino, zaman)
			hata(err)
			http.Redirect(w, r, "/", 302)
		}
		http.Redirect(w, r, r.URL.String(), 302)
	}

	html, err := template.ParseFiles("root"+r.URL.String()+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)
	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.YoneticiAd = profil.Ad
	html.Execute(w, veri)
}

func siteDuzenle(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type SiteDuzenle struct {
		Apartman
		AktifSiteBilgi Site
		SiteSayi       int
	}
	var veri SiteDuzenle

	if r.Method == "POST" {
		if r.FormValue("siteduzenle") != "" {
			sitead := r.FormValue("sitead")
			yoneticiad := r.FormValue("yoneticiad")
			sitefoto := r.FormValue("sitefoto")
			hakkinda := r.FormValue("hakkinda")
			telefon := r.FormValue("telefon")
			ceptelefon := r.FormValue("ceptelefon")
			adres := r.FormValue("adres")
			vergidaire := r.FormValue("vergidaire")
			vergino := r.FormValue("vergino")
			zaman := strconv.FormatInt(time.Now().Unix(), 10)

			_, fotoheader, err := r.FormFile("foto")
			if err == nil {
				sitefoto = fotoheader.Filename
				dosya, err := fotoheader.Open()
				hata(err)
				b, err := ioutil.ReadAll(dosya)
				dosya.Close()
				hata(err)
				dosya2, err := os.Create("root/img/" + sitefoto)
				hata(err)
				dosya2.Write(b)
				dosya2.Close()
			}

			id, err := strconv.Atoi(aktifSiteGetir())
			hata(err)
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			_, err = db.Exec("update site set ad = ?,yoneticiad = ?,foto = ?,hakkinda = ?,telefon = ?,ceptelefon = ?,adres = ?,vergidairesi = ?,vergino = ?,zaman = ? where id = ?",
				sitead, yoneticiad, sitefoto, hakkinda, telefon, ceptelefon, adres, vergidaire, vergino, zaman, id)
			hata(err)
		} else if (r.FormValue("sitesil") != "") {
			siteSil()
		}
		http.Redirect(w, r, r.URL.String(), 302)
	}

	veri.SiteSayi = siteSayiGetir()

	html, err := template.ParseFiles("root"+r.URL.String()+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)
	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.AktifSiteBilgi = aktifSiteBilgi()
	html.Execute(w, veri)
}

func siteSayiGetir() int {
	db := connectDB()
	rows, err := db.Query("select id from site")
	hata(err)
	var sayi = 0
	for rows.Next() {
		sayi += 1
	}
	db.Close()
	return sayi
}

func siteSil() {
	sayi := siteSayiGetir()
	if sayi > 1 {
		id, err := strconv.Atoi(aktifSiteGetir())
		hata(err)
		db := connectDB()
		_, err = db.Exec("delete from tanimlirapor where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from borclandirma where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from odeme where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from dairehareket where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from daire where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from firma where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from gelirgider where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from kisi where siteid = ?", id)
		if err != nil {
			return
		}
		_, err = db.Exec("delete from site where id = ?", id)
		hata(err)

		//Diğer Site Seçimi
		err = db.QueryRow("select id from site order by id asc limit 1").Scan(&id)
		if err != nil {
			fmt.Println("hata:", err)
		}
		hata(err)
		fmt.Println("id", id)
		_, err = db.Exec("update diger set value = ? where key = 'aktifsite'", id)
		hata(err)
	}
}
