package apartmanx

import (
	"html/template"
	"net/http"
	"strconv"
	"database/sql"
)

func anasayfa(w http.ResponseWriter, r *http.Request) {
	type Anasayfa struct {
		Apartman
		AktifSite int
		Siteler   []Site
	}
	var html *template.Template
	var err error
	var veri Anasayfa

	if r.Method == "POST" {
		if r.FormValue("sitesec") != "" {
			id := r.FormValue("sitesec")
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			_, err = db.Exec("update diger set value = ? where key = 'aktifsite'", id)
			hata(err)
		}
		http.Redirect(w, r, r.URL.String(), 302)
	}

	html, err = template.ParseFiles("root/index.html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.Siteler = siteleriGetir()
	veri.AktifSite, err = strconv.Atoi(aktifSiteGetir())
	hata(err)
	veri.SiteAd = aktifSiteAd()
	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = true
	html.Execute(w, veri)

}

func giris(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.FormValue("giris") != "" {
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			var sifre string
			err = db.QueryRow("select sifre from profil").Scan(&sifre)
			hata(err)
			if sifre == r.FormValue("sifre") {
				_, err = db.Exec("update diger set value = 1 where key = 'oturum'")
				hata(err)
			}
			http.Redirect(w, r, "http://localhost:4182", 302)
		}
	}
	html, err := template.ParseFiles("root/giris.html", "root/genel/footer.html")
	if err != nil {
		panic(err)
	}
	html.Execute(w, nil)
}

func cikis(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	_, err = db.Exec("update diger set value = 0 where key = 'oturum'")
	hata(err)
	http.Redirect(w, r, "http://localhost:4182", 302)
}

func iletisim(w http.ResponseWriter, r *http.Request) {
	type Iletisim struct {
		Apartman
	}
	var veri Iletisim
	html, err := template.ParseFiles("root/iletisim.html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)
	veri.SiteAd = aktifSiteAd()
	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	html.Execute(w, veri)
}
