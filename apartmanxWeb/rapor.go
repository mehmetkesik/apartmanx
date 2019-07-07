package apartmanx

import (
	"net/http"
	"strconv"
	"html/template"
	"strings"
	"sort"
	"os/user"
)

func rapor(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type RaporCikti struct {
		Apartman
		LinkVeri         LinkVeri
		Rapor            Rapor
		AramaDurum       bool
		Saf              bool
		RaporBaslik      template.HTML
		RaporBaslikSatir int
		RaporNot         template.HTML
		RaporNotSatir    int
	}
	var veri RaporCikti
	link := strings.Split(r.URL.String(), "?")[0]

	if r.Method == "POST" {
		if r.FormValue("pdfaktar") != "" {
			raporlink := r.FormValue("raporlink")
			if raporlink != "" {
				raporbaslik := r.FormValue("raporbaslik")
				rapornot := r.FormValue("rapornot")
				db := connectDB()
				_, err := db.Exec("update diger set value = ? where key = ?", raporbaslik, "raporbaslik")
				hata(err)
				_, err = db.Exec("update diger set value = ? where key = ?", rapornot, "rapornot")
				hata(err)
				db.Close()

				raporlink += "&saf=xxx"
				kullanici, err := user.Current()
				hata(err)
				html2pdf(raporlink, kullanici.HomeDir+"/Desktop/rapor.pdf")
			}
		}
		http.Redirect(w, r, r.URL.String(), 302)
	}

	if r.URL.Query().Get("rapor") != "" {
		veri.AramaDurum = true
		var err error

		if r.URL.Query().Get("saf") != "" {
			veri.Saf = true
			db := connectDB()
			err := db.QueryRow("select value from diger where key = ?", "raporbaslik").Scan(&veri.RaporBaslik)
			hata(err)
			err = db.QueryRow("select value from diger where key = ?", "rapornot").Scan(&veri.RaporNot)
			hata(err)
			db.Close()
			raporbasliklar := strings.Split(string(veri.RaporBaslik), "\n")
			veri.RaporBaslikSatir = len(raporbasliklar)
			for _, k := range raporbasliklar {
				if len(k) > 102 {
					i := len(k) / 102
					veri.RaporBaslikSatir += i
				}
			}
			veri.RaporBaslikSatir += 1
			rapornotlar := strings.Split(string(veri.RaporNot), "\n")
			veri.RaporNotSatir = len(rapornotlar)
			for _, k := range rapornotlar {
				if len(k) > 57 {
					i := len(k) / 57
					veri.RaporNotSatir += i
				}
			}
			veri.RaporNotSatir += 1
		}

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

		sort.Slice(veri.Rapor.Borclandirma, func(i, j int) bool { //Borçları tarihe göre sıralar.
			return veri.Rapor.Borclandirma[i].Borc.Tarih < veri.Rapor.Borclandirma[j].Borc.Tarih
		})

		sort.Slice(veri.Rapor.Odeme, func(i, j int) bool { //Borçları tarihe göre sıralar.
			return veri.Rapor.Odeme[i].Tarih < veri.Rapor.Odeme[j].Tarih
		})

	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}
