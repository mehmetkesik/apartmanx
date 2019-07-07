package apartmanx

import (
	"net/http"
	"database/sql"
	"html/template"
	"time"
	"strconv"
	"strings"
)

type ProfilX struct {
	Id       int
	Ad       string
	Soyad    string
	Mail     string
	Telefon  string
	Yas      string
	Cinsiyet string
	Adres    string
	Sifre    string
	Zaman    string
}

func profil(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type Profil struct {
		Apartman
		Kisi       ProfilX
		ProfilKayitDurum bool
	}
	var veri Profil
	veri.Kisi.profilGetir()

	link := strings.Split(r.URL.String(), "?")[0]

	if r.Method == "POST" {
		if r.FormValue("profilkaydet") != "" {
			veri.Kisi.Ad = r.FormValue("ad")
			veri.Kisi.Soyad = r.FormValue("soyad")
			veri.Kisi.Mail = r.FormValue("mail")
			veri.Kisi.Telefon = r.FormValue("telefon")
			veri.Kisi.Yas = r.FormValue("yas")
			veri.Kisi.Cinsiyet = r.FormValue("cinsiyet")
			veri.Kisi.Adres = r.FormValue("adres")
			veri.Kisi.Sifre = r.FormValue("sifre")
			veri.Kisi.Zaman = strconv.FormatInt(time.Now().Unix(), 10)
			if veri.Kisi.Sifre != "" && veri.Kisi.Ad != "" && veri.Kisi.Soyad != "" {
				veri.Kisi.profilKaydet()
				http.Redirect(w, r, link+"?durum=1", 302)
			}
		}
		http.Redirect(w, r, link, 302)
	}
	veri.ProfilKayitDurum = false
	if r.URL.Query().Get("durum") == "1" {
		veri.ProfilKayitDurum = true
	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)
	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	html.Execute(w, veri)
}

func (self *ProfilX) profilGetir() {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	defer db.Close()
	hata(err)
	err = db.QueryRow("select * from profil where id = 1").Scan(&self.Id, &self.Ad, &self.Soyad,
		&self.Mail, &self.Telefon, &self.Yas, &self.Cinsiyet, &self.Adres, &self.Sifre, &self.Zaman)
	hata(err)
}

func (self *ProfilX) profilKaydet() {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	defer db.Close()
	hata(err)
	_, err = db.Exec("update profil set ad=?,soyad=?,mail=?,telefon=?,yas=?,cinsiyet=?,adres=?,sifre=?,zaman=?", self.Ad, self.Soyad,
		self.Mail, self.Telefon, self.Yas, self.Cinsiyet, self.Adres, self.Sifre, self.Zaman)
	hata(err)
}
