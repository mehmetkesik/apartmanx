package apartmanx

import (
	"strconv"
	"database/sql"
	"strings"
	"time"
	"fmt"
)

type BorclandirmaRapor struct {
	Borc     Borclandirma
	BorcluId int
	BorcluAd string
	BorcGun  int
	BorcPayi float64
}

func daireBorclari(blok string, daireno int, bastarih, bittarih string) ([]BorclandirmaRapor, float64, float64) {
	var daireBorclariZaman []Borclandirma
	daireBorclariZaman = daireBorclariGetir(blok, daireno, bastarih, bittarih)
	//zamangrup 1 = anlık, zamangrup 2 = aylık
	dairehareketleri := daireHareketlerZamanGetir(blok, daireno, bastarih, bittarih)
	var borclandirmaRaporlar []BorclandirmaRapor
	var toplamBorc, toplamOdeme float64
	for _, borc := range daireBorclariZaman {
		tutborc, err := strconv.ParseFloat(borc.Miktar, 64)
		hata(err)
		toplamBorc += tutborc
		toplamOdeme += borc.OdemeToplam
		if borc.Zamangrup == "1" { //anlık borçlar
			tut := anlikSorumluBul(borc, dairehareketleri)
			if tut != nil {
				borclandirmaRaporlar = append(borclandirmaRaporlar, *tut)
			}
		} else if borc.Zamangrup == "2" { //aylık borçlar
			borclandirmaRaporlar = append(borclandirmaRaporlar, aylikSorumluBul(borc, dairehareketleri)...)
		}
	}
	return borclandirmaRaporlar, toplamBorc, toplamOdeme
}

func kisiBorclari(kisiid int, blok string, daireno int, bastarih, bittarih string) {

}

func daireBorclariGetir(blok string, daireno int, bastarih, bittarih string) []Borclandirma {
	var bastarihaylik, bittarihaylik string
	if bastarih != "" {
		bastarihdizi := strings.Split(bastarih, "-")
		bastarihaylik = bastarihdizi[0] + "-" + bastarihdizi[1]
	}
	if bittarih != "" {
		bittarihdizi := strings.Split(bittarih, "-")
		bittarihaylik = bittarihdizi[0] + "-" + bittarihdizi[1]
	}

	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var borclandirmalar []Borclandirma
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	var rows *sql.Rows
	if bastarih == "" && bittarih == "" { //tarih yok
		rows, err = db.Query("select * from borclandirma where siteid = ? and blok = ? and daireno = ? order by tarih asc", siteid, blok, daireno)
		hata(err)
	} else if bastarih != "" && bittarih == "" { //başlama tarihi var
		rows, err = db.Query("select * from borclandirma where siteid = ? and blok = ? and daireno = ? and (tarih >= ? or tarih >= ?) order by tarih asc",
			siteid, blok, daireno, bastarih, bastarihaylik)
		hata(err)
	} else if bastarih == "" && bittarih != "" { //bitiş tarihi var
		rows, err = db.Query("select * from borclandirma where siteid = ? and blok = ? and daireno = ? and (tarih <= ? or tarih <= ?) order by tarih asc",
			siteid, blok, daireno, bittarih, bittarihaylik)
		hata(err)
	} else { //iki tarihte var
		rows, err = db.Query("select * from borclandirma where siteid = ? and blok = ? and daireno = ? and (tarih >= ? or tarih >= ?) and (tarih <= ? or tarih <= ?) order by tarih asc",
			siteid, blok, daireno, bastarih, bastarihaylik, bittarih, bittarihaylik)
		hata(err)
	}

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

func daireHareketlerZamanGetir(blok string, daireno int, bastarih, bittarih string) []DaireHareket {
	var dairehareketler []DaireHareket
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	id, err := strconv.Atoi(aktifSiteGetir())
	hata(err)

	var rows *sql.Rows
	if bastarih == "" && bittarih == "" { //tarih yok
		rows, err = db.Query("select * from dairehareket where siteid = ? and blok = ? and daireno = ?", id, blok, daireno)
		hata(err)
	} else if bastarih != "" && bittarih == "" { //başlama tarihi var
		rows, err = db.Query("select * from dairehareket where siteid = ? and blok = ? and daireno = ? and tarih >= ?",
			id, blok, daireno, bastarih)
		hata(err)
	} else if bastarih == "" && bittarih != "" { //bitiş tarihi var
		rows, err = db.Query("select * from dairehareket where siteid = ? and blok = ? and daireno = ? and tarih <= ?",
			id, blok, daireno, bittarih)
		hata(err)
	} else { //iki tarihte var
		rows, err = db.Query("select * from dairehareket where siteid = ? and blok = ? and daireno = ? and tarih >= ? and tarih <= ?",
			id, blok, daireno, bastarih, bittarih)
		hata(err)
	}

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

func anlikSorumluBul(borc Borclandirma, dairehareketleri []DaireHareket) *BorclandirmaRapor {
	var onceki DaireHareket
	var borclandirmaRapor BorclandirmaRapor
	for _, hareket := range dairehareketleri { //Eğer borç daire hareket tarihinden önce eklenmişse boş veri ekler. isim "" olur.
		if hareket.Tarih < borc.Tarih {
			onceki = hareket
		} else if hareket.Tarih == borc.Tarih {
			onceki = hareket
			break
		} else {
			break
		}
	}
	borclandirmaRapor.Borc = borc
	if borc.Kime == "1" { //Oturan
		if onceki.KiraciId == -1 {
			borclandirmaRapor.BorcluId = onceki.MulksahibiId
			borclandirmaRapor.BorcluAd = onceki.Mulksahibi
		} else {
			borclandirmaRapor.BorcluId = onceki.KiraciId
			borclandirmaRapor.BorcluAd = onceki.Kiraci
		}
	} else if borc.Kime == "2" { //Mülk Sahibi
		borclandirmaRapor.BorcluId = onceki.MulksahibiId
		borclandirmaRapor.BorcluAd = onceki.Mulksahibi
	} else if borc.Kime == "3" { //Kiracı
		borclandirmaRapor.BorcluId = onceki.KiraciId
		borclandirmaRapor.BorcluAd = onceki.Kiraci
	}
	if borclandirmaRapor.BorcluId != -1 {
		return &borclandirmaRapor
	}
	return nil
}

func aylikSorumluBul(borc Borclandirma, dairehareketleri []DaireHareket) []BorclandirmaRapor {
	var hareketler []DaireHareket
	var borclandirmaRapor []BorclandirmaRapor

	//Burası tarihin ay kısmını bir arttırmaktadır.
	//borclandirmaRapor.Borc = borc
	var borcust string
	borcustdizi := strings.Split(borc.Tarih, "-")
	tut, err := strconv.Atoi(borcustdizi[1])
	hata(err)
	if tut < 12 {
		tut += 1
		borcust = borcustdizi[0] + "-" + strconv.Itoa(tut)
	} else {
		tut2, err := strconv.Atoi(borcustdizi[0])
		hata(err)
		tut2 += 1
		borcust = strconv.Itoa(tut2) + "-01"
	}
	//Ay kısmını bir arttırma bitti.

	borcalt := borc.Tarih + "-01"
	borcust = borcust + "-01"

	var once = false
	var orta = false

	var gunler []int
	var onceki DaireHareket
	gunler = append(gunler, 1)
	for _, hareket := range dairehareketleri {
		if hareket.Tarih < borcalt {
			onceki = hareket
			once = true
		} else if hareket.Tarih >= borcalt && hareket.Tarih < borcust {
			orta = true
			gun, err := strconv.Atoi(strings.Split(hareket.Tarih, "-")[2])
			hata(err)
			xyz := gun - gunler[len(gunler)-1]
			gunler[len(gunler)-1] = xyz
			gunler = append(gunler, gun)
			hareketler = append(hareketler, hareket)
		} else {
			break
		}
	}
	var aycekim = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	aysayi, err := strconv.Atoi(strings.Split(borc.Tarih, "-")[1])
	hata(err)
	aycekimsayi := aycekim[aysayi-1]
	xyz := (aycekimsayi + 1) - gunler[len(gunler)-1] //Burada ayın kaç çektiğine göre işlem yapılacak.
	gunler[len(gunler)-1] = xyz

	if !once && !orta {
		return borclandirmaRapor
	} else if once && !orta {
		var br BorclandirmaRapor
		brtut := kimBul(borc, br, onceki)
		brtut.BorcGun = aycekimsayi
		miktar, err := strconv.ParseFloat(brtut.Borc.Miktar, 64)
		hata(err)
		brtut.BorcPayi = miktar
		borclandirmaRapor = append(borclandirmaRapor, *brtut)
		return borclandirmaRapor
	} else {
		if once {
			var brx BorclandirmaRapor
			tutbr := kimBul(borc, brx, onceki)
			if tutbr != nil {
				tutbr.BorcGun = gunler[0]
				miktar, err := strconv.ParseFloat(tutbr.Borc.Miktar, 64)
				hata(err)
				miktarx, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(tutbr.BorcGun)/float64(aycekimsayi))*miktar), 64)
				hata(err)
				tutbr.BorcPayi = miktarx
				if miktarx != 0 {
					borclandirmaRapor = append(borclandirmaRapor, *tutbr)
				}
			}
		}
		for i, hareket := range hareketler {
			var br BorclandirmaRapor
			tutbr := kimBul(borc, br, hareket)
			if tutbr != nil {
				tutbr.BorcGun = gunler[i+1]
				miktar, err := strconv.ParseFloat(tutbr.Borc.Miktar, 64)
				hata(err)
				miktarx, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(tutbr.BorcGun)/float64(aycekimsayi))*miktar), 64)
				hata(err)
				tutbr.BorcPayi = miktarx
				borclandirmaRapor = append(borclandirmaRapor, *tutbr)
			}
		}
	}

	var brmap = make(map[int]BorclandirmaRapor)
	for _, br := range borclandirmaRapor {
		x := brmap[br.BorcluId]
		x.Borc = br.Borc
		x.BorcluId = br.BorcluId
		x.BorcPayi = br.BorcPayi
		x.BorcluAd = br.BorcluAd
		x.BorcGun += br.BorcGun
		brmap[br.BorcluId] = x
	}
	var brapor []BorclandirmaRapor
	for _, rapor := range brmap {
		brapor = append(brapor, rapor)
	}
	return brapor
}

func kimBul(borc Borclandirma, br BorclandirmaRapor, hareket DaireHareket) *BorclandirmaRapor {
	if borc.Kime == "1" { //Oturan
		if hareket.KiraciId == -1 {
			br.BorcluId = hareket.MulksahibiId
			br.BorcluAd = hareket.Mulksahibi
		} else {
			br.BorcluId = hareket.KiraciId
			br.BorcluAd = hareket.Kiraci
		}
	} else if borc.Kime == "2" { //Mülk Sahibi
		br.BorcluId = hareket.MulksahibiId
		br.BorcluAd = hareket.Mulksahibi
	} else if borc.Kime == "3" { //Kiracı
		br.BorcluId = hareket.KiraciId
		br.BorcluAd = hareket.Kiraci
	}
	if br.BorcluId != -1 {
		br.Borc = borc
		return &br
	}
	return nil
}
