package apartmanx

import (
	"net/http"
	"text/template"
	"encoding/json"
	"fmt"
	"strings"
	"database/sql"
	"strconv"
	"time"
	"sort"
)

type Tipi struct {
	Id int
	Ad string
}

type Daire struct {
	Id           int
	Siteid       int
	Blok         string
	Daireno      string
	Mulksahibi   string
	MulksahibiId int
	Kiraci       string
	KiraciId     int
	TipiId       int
	Zaman        string
	Ekstra       map[string][]string
}

type DaireHareket struct {
	Id           int
	Siteid       int
	Blok         string
	Daireno      string
	Mulksahibi   string
	MulksahibiId int
	Kiraci       string
	KiraciId     int
	Degisen      [2]bool
	Tarih        string
	Zaman        string
}

func daireKayitDuzenleme(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type DaireKayitIslemi struct {
		Apartman
		DaireKayitDurum bool
		DaireSilDurum   bool
		DaireBilgi      Daire
		DaireAra        bool
		EkstraVarMi     bool
		Link            string
		DaireBlok       []string
		KisilerBilgi    []Kisi
		Tipi            []Tipi
	}
	var veri DaireKayitIslemi
	veri.DaireBlok = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U", "V", "Y", "Z"}
	veri.Tipi = tipiGetir()

	link := strings.Split(r.URL.String(), "?")[0]
	veri.KisilerBilgi = kisilerGetir()

	if r.Method == "POST" {
		if r.FormValue("dairekayit") != "" {
			var tutx = make(map[string]string)
			var tut = make(map[string][]string)
			var say = 1
			for key, value := range r.Form {
				if strings.HasPrefix(key, "ekstra") {
					tutstr := strings.Split(key, ":")
					xyz := value[0]      //alanın değeri
					value[0] = tutstr[1] //tutstr[1] alanın adı, tutstr[0] name'dir.
					value = append(value, xyz)
					tut["ekstra"+strconv.Itoa(say)] = value
					say += 1
				} else {
					tutx[key] = value[0]
				}
			}
			if tutx["kiraciId"] == "" {
				tutx["kiraciId"] = "-1"
			}
			b, err := json.Marshal(tut)
			if err != nil {
				fmt.Println(err)
			}
			siteid, err := strconv.Atoi(aktifSiteGetir())
			hata(err)
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			db, err := sql.Open("sqlite3", "apartmanx.db")
			defer db.Close()
			hata(err)

			var mulksahibi, kiraci int
			var degisim = false
			var degisen = ""
			err = db.QueryRow("select mulksahibi,kiraci from daire where blok = ? and daireno = ? and siteid = ?",
				tutx["blok"], tutx["daireno"], siteid).Scan(&mulksahibi, &kiraci)
			if err != nil {
				degisim = true
				degisen = "-1"
				_, err := db.Exec("insert into daire (siteid,blok,daireno,mulksahibi,kiraci,tipi,zaman,ekstra) values "+
					"(?,?,?,?,?,?,?,?)", siteid, tutx["blok"], tutx["daireno"], tutx["mulksahibiId"], tutx["kiraciId"],
					tutx["tipiId"], zaman, string(b))
				hata(err)
			} else {
				if tutx["mulksahibiId"] != strconv.Itoa(mulksahibi) {
					degisim = true
					degisen = "0" //Ev Sahibinin Değiştiğini Gösterir
				}
				if tutx["kiraci"] != strconv.Itoa(kiraci) {
					degisim = true
					if degisen == "" {
						degisen += "1"
					} else {
						degisen += ",1"
					}
				}
				_, err := db.Exec("update daire set blok=?,daireno=?,mulksahibi=?,kiraci=?,tipi=?,zaman=?,ekstra=?"+
					" where daireno = ? and blok = ? and siteid = ?",
					tutx["blok"], tutx["daireno"], tutx["mulksahibiId"], tutx["kiraciId"], tutx["tipiId"], zaman,
					string(b), tutx["daireno"], tutx["blok"], siteid)
				hata(err)
			}

			if degisim {
				var hareketid int
				zamanint, err := strconv.ParseInt(zaman, 10, 64)
				hata(err)
				tarih := strings.Split(time.Unix(zamanint, 0).String(), " ")[0]
				err = db.QueryRow("select id from dairehareket where siteid = ? and blok = ? and daireno = ? and tarih = ?",
					siteid, tutx["blok"], tutx["daireno"], tarih).Scan(&hareketid)
				if err != nil {
					fmt.Println("insrt amk")
					_, err := db.Exec("insert into dairehareket (siteid,blok,daireno,mulksahibi,kiraci,degisen,tarih,zaman) values "+
						"(?,?,?,?,?,?,?,?)", siteid, tutx["blok"], tutx["daireno"], tutx["mulksahibiId"], tutx["kiraciId"], degisen,
						tarih, zaman)
					hata(err)
				} else {
					fmt.Println("update amk")
					_, err := db.Exec("update dairehareket set siteid = ?,blok = ?,daireno = ?,mulksahibi = ?,kiraci = ?,"+
						"degisen = ?,tarih = ?,zaman = ? where id = ?", siteid, tutx["blok"], tutx["daireno"], tutx["mulksahibiId"],
						tutx["kiraciId"], degisen, tarih, zaman, hareketid)
					hata(err)
				}

			}

			http.Redirect(w, r, link+"?durum=1", 302)
		} else if r.FormValue("dairesil") != "" {
			if r.FormValue("blok") != "" && r.FormValue("daireno") != "" {
				siteid, err := strconv.Atoi(aktifSiteGetir())
				hata(err)
				blok := r.FormValue("blok")
				dairenostr := r.FormValue("daireno")
				daireno, err := strconv.Atoi(dairenostr)
				hata(err)
				db := connectDB()
				_, err = db.Exec("delete from daire where siteid = ? and blok = ? and daireno = ?", siteid, blok, daireno)
				hata(err)
				_, err = db.Exec("delete from dairehareket where siteid = ? and blok = ? and daireno = ?", siteid, blok, daireno)
				hata(err)
				var borcidler []int
				rows, err := db.Query("select id from borclandirma where  siteid = ? and blok = ? and daireno = ?", siteid, blok, daireno)
				for rows.Next() {
					var borcid int
					err = rows.Scan(&borcid)
					hata(err)
					borcidler = append(borcidler, borcid)
				}
				_, err = db.Exec("delete from borclandirma where siteid = ? and blok = ? and daireno = ?", siteid, blok, daireno)
				hata(err)
				for _, k := range borcidler {
					_, err = db.Exec("delete from odeme where siteid = ? and borcid = ?", siteid, k)
					hata(err)
				}
			}
			http.Redirect(w, r, link+"?durum=2", 302)
		}

		http.Redirect(w, r, link, 302)
	}

	siteid, err := strconv.Atoi(aktifSiteGetir())
	if r.URL.Query().Get("daireno") != "" && r.URL.Query().Get("blok") != "" {
		var daire Daire
		var ekstra string
		daireno := r.URL.Query().Get("daireno")
		blok := r.URL.Query().Get("blok")
		db, err := sql.Open("sqlite3", "apartmanx.db")
		hata(err)
		err = db.QueryRow("select * from daire where daireno = ? and blok = ? and siteid = ?", daireno, blok, siteid).
			Scan(&daire.Id, &daire.Siteid, &daire.Blok, &daire.Daireno, &daire.MulksahibiId, &daire.KiraciId, &daire.TipiId, &daire.Zaman, &ekstra)

		var tutsoyad string
		db.QueryRow("select ad,soyad from kisi where id = ?", daire.MulksahibiId).Scan(&daire.Mulksahibi, &tutsoyad)
		daire.Mulksahibi += " " + tutsoyad
		if daire.KiraciId != -1 {
			db.QueryRow("select ad,soyad from kisi where id = ?", daire.KiraciId).Scan(&daire.Kiraci, &tutsoyad)
			daire.Kiraci += " " + tutsoyad
		}

		if err != nil {
			http.Redirect(w, r, link, 302)
		} else {
			err = json.Unmarshal([]byte(ekstra), &daire.Ekstra)
			hata(err)
			veri.DaireBilgi = daire
			veri.DaireAra = true
			if len(daire.Ekstra) > 0 {
				veri.EkstraVarMi = true
			} else {
				veri.EkstraVarMi = false
			}
		}
	}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.DaireKayitDurum = false
	if r.URL.Query().Get("durum") != "" {
		durum := r.URL.Query().Get("durum")
		if durum == "1" {
			veri.DaireKayitDurum = true
		} else if durum == "2" {
			veri.DaireSilDurum = true
		}
	}

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	veri.Link = link
	html.Execute(w, veri)
}

//////////////////////////////////////////////////////////////
func daireListesi(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type DaireKayitListesi struct {
		Apartman
		Daireler []Daire
		Link     string
		Tipi     []Tipi
	}
	var veri DaireKayitListesi
	link := strings.Split(r.URL.String(), "?")[0]

	veri.Tipi = tipiGetir()

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.Link = link
	veri.Daireler = dairelerGetir()
	sort.Slice(veri.Daireler, func(i, j int) bool { //Daireleri önce bloğa sonra dairenosuna göre sıralar.
		if !(veri.Daireler[i].Blok < veri.Daireler[j].Blok) {
			if veri.Daireler[i].Blok == veri.Daireler[j].Blok {
				return veri.Daireler[i].Daireno < veri.Daireler[j].Daireno
			} else {
				return veri.Daireler[i].Blok < veri.Daireler[j].Blok
			}
		}
		return veri.Daireler[i].Blok < veri.Daireler[j].Blok
	})
	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}

//////////////////////////////////////////////////////////////
func daireHareketleri(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	type DaireKayitHareketleri struct {
		Apartman
		DaireHareketler []DaireHareket
		DaireArama      bool
		Link            string
		DaireBlok       []string
		Tipi            []Tipi
	}
	var veri DaireKayitHareketleri
	link := strings.Split(r.URL.String(), "?")[0]

	veri.Tipi = tipiGetir()

	veri.DaireBlok = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U", "V", "Y", "Z"}

	html, err := template.ParseFiles("root"+link+".html", "root/genel/header.html", "root/genel/footer.html")
	hata(err)

	veri.Link = link
	veri.DaireArama = false
	if r.URL.Query().Get("blok") != "" && r.URL.Query().Get("daireno") != "" {
		veri.DaireHareketler = daireHareketlerGetir(r.URL.Query().Get("blok"), r.URL.Query().Get("daireno"))
		veri.DaireArama = true
	}

	veri.HeaderMenu = getMenu()
	veri.SiteEkleDurum = false
	veri.SiteAd = aktifSiteAd()
	veri.Ad = adBul(link, veri.HeaderMenu)
	html.Execute(w, veri)
}
