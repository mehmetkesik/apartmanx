package apartmanx

import (
	"net/http"
	"database/sql"
	"strconv"
	"time"
	"fmt"
	"strings"
	"encoding/json"
)

func post(w http.ResponseWriter, r *http.Request) {
	if !oturumKontrol(w, r) {
		var j = make(map[string]string)
		j["sonuc"] = "Oturum Başarısız!"
		sonuc, _ := json.Marshal(j)
		fmt.Fprint(w, string(sonuc))
		return
	}
	if r.Method == "POST" {
		siteid, err := strconv.Atoi(aktifSiteGetir())
		hata(err)

		//// Kişi Güncelleme
		if r.FormValue("kisiguncelle") != "" {
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			_, err = db.Exec("update kisi  set ad = ?,soyad = ?,mail = ?,telefon = ?,cinsiyet = ?,yas = ?,zaman = ? where id = ? and siteid = ?",
				r.FormValue("ad"), r.FormValue("soyad"), r.FormValue("mail"), r.FormValue("telefon"),
				r.FormValue("cinsiyet"), r.FormValue("yas"), zaman, r.FormValue("id"), siteid)
			hata(err)

			sdizi := strings.Split(time.Unix(time.Now().Unix(), 0).String(), " ")
			zaman = sdizi[0] + " " + sdizi[1]
			var j = make(map[string]string)
			j["sonuc"] = "olumlu"
			j["zaman"] = zaman
			sonuc, _ := json.Marshal(j)

			fmt.Fprint(w, string(sonuc))
		}

		//// Firma Güncelleme
		if r.FormValue("firmaguncelle") != "" {
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			_, err = db.Exec("update firma set ad = ?,sahibad = ?,mail = ?,telefon = ?,adres = ?,zaman = ? where id = ? and siteid = ?",
				r.FormValue("ad"), r.FormValue("sahibad"), r.FormValue("mail"), r.FormValue("telefon"),
				r.FormValue("adres"), zaman, r.FormValue("id"), siteid)
			hata(err)

			sdizi := strings.Split(time.Unix(time.Now().Unix(), 0).String(), " ")
			zaman = sdizi[0] + " " + sdizi[1]
			var j = make(map[string]string)
			j["sonuc"] = "olumlu"
			j["zaman"] = zaman
			sonuc, _ := json.Marshal(j)
			fmt.Fprint(w, string(sonuc))
		}

		//// Gelir-Gider Güncelleme
		if r.FormValue("gelirgiderguncelle") != "" {
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			zaman := strconv.FormatInt(time.Now().Unix(), 10)
			_, err = db.Exec("update gelirgider set ad = ?,zaman = ? where id = ? and siteid = ?",
				r.FormValue("ad"), zaman, r.FormValue("id"), siteid)
			hata(err)

			sdizi := strings.Split(time.Unix(time.Now().Unix(), 0).String(), " ")
			zaman = sdizi[0] + " " + sdizi[1]
			var j = make(map[string]string)
			j["sonuc"] = "olumlu"
			j["zaman"] = zaman
			sonuc, _ := json.Marshal(j)

			fmt.Fprint(w, string(sonuc))
		}

		if r.FormValue("odemeyap") != "" {
			var durum = true
			var tutx = make(map[string]string)
			for key, value := range r.Form {
				tutx[key] = value[0]
				if value[0] == "" {
					durum = false
				}
			}

			if durum {
				sonid := odemeyaz(tutx)
				var j = make(map[string]string)
				j["sonuc"] = "olumlu"
				j["id"] = strconv.FormatInt(sonid, 10)
				sonuc, _ := json.Marshal(j)
				fmt.Fprint(w, string(sonuc))
			} else {
				var j = make(map[string]string)
				j["sonuc"] = "olumsuz"
				sonuc, _ := json.Marshal(j)
				fmt.Fprint(w, string(sonuc))
			}
		}

		if r.FormValue("borcsil") != "" {
			borcidstr := r.FormValue("id")
			borcid, err := strconv.Atoi(borcidstr)
			hata(err)
			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()
			_, err = db.Exec("delete from borclandirma where id = ? and siteid = ?", borcid, siteid)
			hata(err)
			_, err = db.Exec("delete from odeme where borcid = ? and siteid = ?", borcid, siteid)
			hata(err)
			var j = make(map[string]string)
			j["sonuc"] = "olumlu"
			sonuc, _ := json.Marshal(j)
			fmt.Fprint(w, string(sonuc))
		}

		if r.FormValue("odemesil") != "" {
			borcidstr := r.FormValue("borcid")
			odemeidstr := r.FormValue("odemeid")
			borcid, err := strconv.Atoi(borcidstr)
			hata(err)
			odemeid, err := strconv.Atoi(odemeidstr)
			hata(err)

			db, err := sql.Open("sqlite3", "apartmanx.db")
			hata(err)
			defer db.Close()

			_, err = db.Exec("delete from odeme where id = ? and borcid = ? and siteid = ?", odemeid, borcid, siteid)
			hata(err)

			var j = make(map[string]string)
			j["sonuc"] = "olumlu"
			sonuc, _ := json.Marshal(j)
			fmt.Fprint(w, string(sonuc))
		}

		if r.FormValue("sablonkayit") != "" {
			sablonad := r.FormValue("ad")
			sablonlink := r.FormValue("link")
			if sablonad == "" || sablonlink == "" {
				var j = make(map[string]string)
				j["sonuc"] = "olumsuz"
				sonuc, _ := json.Marshal(j)
				fmt.Fprint(w, string(sonuc))
			} else {
				zaman := strconv.FormatInt(time.Now().Unix(), 10)
				db := connectDB()

				rows, err := db.Query("select count(*) as sayi from tanimlirapor where siteid = ?", siteid)
				hata(err)
				var sayi int
				for rows.Next() {
					err := rows.Scan(&sayi)
					hata(err)
				}
				if sayi >= 500 {
					var j = make(map[string]string)
					j["sonuc"] = "sinir"
					sonuc, _ := json.Marshal(j)
					fmt.Fprint(w, string(sonuc))
				} else {
					_, err = db.Exec("insert into tanimlirapor (siteid,ad,link,zaman) values (?,?,?,?)", siteid, sablonad, sablonlink, zaman)
					hata(err)

					var j = make(map[string]string)
					j["sonuc"] = "olumlu"
					sonuc, _ := json.Marshal(j)
					fmt.Fprint(w, string(sonuc))
				}
			}
		}

		if r.FormValue("sablonsil") != "" {
			sablonsilidstr := r.FormValue("id")
			sablonsilid, err := strconv.Atoi(sablonsilidstr)
			hata(err)
			db := connectDB()
			_, err = db.Exec("delete from tanimlirapor where id = ? and siteid = ?", sablonsilid, siteid)
			hata(err)
			var j = make(map[string]string)
			j["sonuc"] = "olumlu"
			sonuc, _ := json.Marshal(j)
			fmt.Fprint(w, string(sonuc))
		}

	} //post son
}
