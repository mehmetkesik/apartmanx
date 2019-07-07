package apartmanx

import (
	"net/http"
	"html/template"
	"strings"
	"errors"
)

type Odeme struct {
	Id     int
	Siteid int
	Borcid int
	Kisiid int
	Miktar string
	Tarih  string
	Zaman  string
	Kisiad string // Bu veritabanında yok
}

type Borclandirma struct {
	Id           int
	Siteid       int
	Blok         string
	Daireno      string
	Miktar       string
	Kime         string
	Zamangrup    string
	Tarih        string
	Gelirgiderid int
	Gelirgiderad string
	Firmaid      int
	Firmaad      string
	Not          string
	Zaman        string
	Odemeler     []Odeme //Bu veritabanında yok
	OdemeToplam  float64
}

func borclandirma(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type BorclandirmaBilgi struct {
		Apartman
		HeaderMenu     []Menu
		Link           string
		Ad             string
		SiteEkleDurum  bool
		SiteAd         string
		BorcKayitDurum int
		DaireBlok      []string
		Tipi           []Tipi
		GelirGiderler  []GelirGider
		Firmalar       []Firma
	}
	var veri BorclandirmaBilgi
	link := strings.Split(r.URL.String(), "?")[0]

	veri.BorcKayitDurum = 0
	veri.DaireBlok = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U", "V", "Y", "Z"}
	veri.Tipi = tipiGetir()
	veri.GelirGiderler = gelirgiderlerGetir()
	veri.Firmalar = firmalarGetir()

	if r.Method == "POST" {
		if r.FormValue("borckayit") != "" {
			var tutx = make(map[string]string)
			for name, value := range r.Form {
				tutx[name] = value[0]
			}
			//dairesecim-1=tekdaire,2=tümdaireler
			//zamangrup-1=anlık,2=aylık
			//kime-1=oturan,2=mülksahibi,3=kiracı
			if tutx["dairesecim"] == "1" {
				if dairekontrol(tutx["blok"], tutx["daireno"]) {
					borcyaz(tutx)
				} else {
					http.Redirect(w, r, link+"?durum=2", 302)
				}
			} else if tutx["dairesecim"] == "2" {
				var daireler = dairelerGetir()
				for _, daire := range daireler {
					tutx["blok"] = daire.Blok
					tutx["daireno"] = daire.Daireno
					borcyaz(tutx)
				}
			} else {
				hata(errors.New("Daire seçim hatası!"))
			}
			http.Redirect(w, r, link+"?durum=1", 302)
		}
		http.Redirect(w, r, link, 302)
	}

	if r.URL.Query().Get("durum") == "1" {
		veri.BorcKayitDurum = 1
	} else if r.URL.Query().Get("durum") == "2" {
		veri.BorcKayitDurum = 2
	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	veri.Link = link
	html.Execute(w, veri)
}

///////////////////////////////////////////////////////////////////////
func odeme(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type OdemeBilgi struct {
		Apartman
		HeaderMenu    []Menu
		Link          string
		Ad            string
		SiteEkleDurum bool
		SiteAd        string
		DaireBlok     []string
		KisilerBilgi  []Kisi
		Borclar       []Borclandirma
		OdemeAra      bool
	}
	var veri OdemeBilgi
	link := strings.Split(r.URL.String(), "?")[0]

	veri.DaireBlok = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U", "V", "Y", "Z"}
	veri.KisilerBilgi = kisilerGetir()

	veri.OdemeAra = false
	if r.URL.Query().Get("borcara") == "tumborclar" {
		veri.Borclar = borclarGetir(true, "", "")
		veri.OdemeAra = true
	} else if r.URL.Query().Get("blok") != "" && r.URL.Query().Get("daireno") != "" {
		veri.Borclar = borclarGetir(false, r.URL.Query().Get("blok"), r.URL.Query().Get("daireno"))
		veri.OdemeAra = true
	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	veri.Link = link
	html.Execute(w, veri)
}
