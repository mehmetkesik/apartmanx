package apartmanx

import (
	"net/http"
	"strings"
	"html/template"
	"io/ioutil"
	"os"
	"gopkg.in/gomail.v2"
	"strconv"
)

func mailGonderme(w http.ResponseWriter, r *http.Request) {
	type MailGonderme struct {
		Apartman
		Link string
	}
	var veri MailGonderme
	link := strings.Split(r.URL.String(), "?")[0]

	if r.Method == "POST" {
		if r.FormValue("mailgonder") != "" {
			var tutx = make(map[string]string)
			for name, value := range r.Form {
				tutx[name] = value[0]
			}
			tutx["mailadres"] = tutx["mailadres"] + "@gmail.com"
			_, dosyaheader, err := r.FormFile("dosya")
			var dosyavarmi bool
			if err == nil {
				dosyavarmi = true
				dosyaname := dosyaheader.Filename
				dosya, err := dosyaheader.Open()
				hata(err)
				b, err := ioutil.ReadAll(dosya)
				hata(err)
				dosya.Close()
				yol := "root/dosya/" + dosyaname
				dosyax, err := os.Create(yol)
				hata(err)
				dosyax.Write(b)
				dosyax.Close()
				tutx["dosya"] = yol
			}
			if tutx["tumkisiler"] == "" {
				mailGonder(tutx)
			} else {
				maillerGonder(tutx)
			}
			if dosyavarmi {
				os.Remove(tutx["dosya"])
			}
		}
		http.Redirect(w, r, link, 302)
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

func mailGonder(tutx map[string]string) {
	if tutx["mailadres"] != "" && tutx["sifre"] != "" && tutx["gidecekmail"] != "" {
		m := gomail.NewMessage()
		m.SetHeader("From", tutx["mailadres"])
		m.SetHeader("To", tutx["gidecekmail"])
		m.SetHeader("Subject", tutx["mailbaslik"])
		m.SetBody("text/html", tutx["mailtext"])
		if tutx["dosya"] != "" {
			m.Attach(tutx["dosya"])
		}

		d := gomail.NewDialer("smtp.gmail.com", 587, tutx["mailadres"], tutx["sifre"])

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			//panic(err)
		}
	}
}

func maillerGonder(tutx map[string]string) {
	if tutx["mailadres"] != "" && tutx["sifre"] != "" {
		siteid, err := strconv.Atoi(aktifSiteGetir())
		hata(err)
		db := connectDB()
		rows, err := db.Query("select mail from kisi where siteid = ?", siteid)
		hata(err)
		var mailler []string
		for rows.Next() {
			var mail string
			rows.Scan(&mail)
			if mail != "" {
				mailler = append(mailler, mail)
			}

		}

		if len(mailler) > 0 {
			m := gomail.NewMessage()
			m.SetHeader("From", tutx["mailadres"])
			m.SetHeader("To", mailler...)
			m.SetHeader("Subject", tutx["mailbaslik"])
			m.SetBody("text/html", tutx["mailtext"])
			if tutx["dosya"] != "" {
				m.Attach(tutx["dosya"])
			}

			d := gomail.NewDialer("smtp.gmail.com", 587, tutx["mailadres"], tutx["sifre"])

			// Send the email to Bob, Cora and Dan.
			if err := d.DialAndSend(m); err != nil {
				//panic(err)
			}
		}
	}
}
