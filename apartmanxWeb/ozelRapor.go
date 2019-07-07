package apartmanx

import (
	"strconv"
	"database/sql"
	"strings"
	"time"
	"fmt"
)

type Rapor struct {
	Borclandirma []OzelRapor
	Odeme        []Odeme
	BorcToplam   float64
	OdemeToplam  float64
	KalanBorc        float64
}

type OzelRapor struct {
	Borc     Borclandirma
	Kisi     Kisi
	BorcGun  int
	BorcPayi float64
}

func ozelRapor(linkveri LinkVeri) Rapor {
	var rapor Rapor
	borclandirma, odemeler, borctoplam, odemetoplam := ozelRaporBorcGetir(linkveri)
	rapor.BorcToplam = borctoplam
	rapor.OdemeToplam = odemetoplam
	rapor.KalanBorc = borctoplam - odemetoplam

	if linkveri.KisiDagit {
		rapor.Borclandirma = kisilereDagit(borclandirma, linkveri.KisiId)
		if linkveri.KisiId != -1 {
			borctoplam = 0
			odemetoplam = 0
			for _, k := range rapor.Borclandirma {
				tutborc, err := strconv.ParseFloat(k.Borc.Miktar, 64)
				hata(err)
				borctoplam += tutborc //Burada Ödemeler Toplanıyor
				for _, k2 := range odemeler {
					if k.Borc.Id == k2.Borcid {
						rapor.Odeme = append(rapor.Odeme, k2)

						tutodeme, err := strconv.ParseFloat(k2.Miktar, 64)
						hata(err)
						odemetoplam += tutodeme //Burada Ödemeler Toplanıyor
					}
				}
			}
			rapor.BorcToplam = borctoplam
			rapor.OdemeToplam = odemetoplam
			rapor.KalanBorc = borctoplam - odemetoplam
		} else {
			rapor.Odeme = odemeler
		}
		return rapor
	} else { //Burası kişilere dağıtılmadan verinin gönderilmesi
		rapor.Odeme = odemeler
		var ozelrapor []OzelRapor
		for _, j := range borclandirma {
			var o OzelRapor
			o.Borc = j
			ozelrapor = append(ozelrapor, o)
		}
		rapor.Borclandirma = ozelrapor
		return rapor
	}
}

//Burası kişilere dağıtılmadığında çalışan fonksiyon
func ozelRaporBorcGetir(linkveri LinkVeri) ([]Borclandirma, []Odeme, float64, float64) {
	var borctoplam, odemetoplam float64
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var sorgustr string
	sorgustr = "select * from borclandirma where siteid = " + strconv.Itoa(siteid)
	if linkveri.Blok != "" {
		sorgustr += " and blok = '" + linkveri.Blok + "'"
	}
	if linkveri.DaireNo != -1 {
		sorgustr += " and daireno = " + strconv.Itoa(linkveri.DaireNo)
	}
	if linkveri.Kime != -1 {
		sorgustr += " and kime = " + strconv.Itoa(linkveri.Kime)
	}
	if linkveri.FirmaId != -1 {
		sorgustr += " and firmaid = " + strconv.Itoa(linkveri.FirmaId)
	}
	if linkveri.GelirGiderId != -1 {
		sorgustr += " and gelirgiderid = " + strconv.Itoa(linkveri.GelirGiderId)
	}
	if linkveri.BasTarih != "" {
		var tut string
		baszamanx := strings.Split(linkveri.BasTarih, "-")
		if baszamanx[2] == "01" {
			ayint, err := strconv.Atoi(baszamanx[1])
			hata(err)
			if ayint > 1 {
				ayint -= 1
				baszamanx[1] = strconv.Itoa(ayint)
				baszamanx[2] = "32" //Burası ayın günlerinden büyük olması için yazıldı.
			} else {
				ayint = 12
				yilint, err := strconv.Atoi(baszamanx[1])
				hata(err)
				yilint -= 1
				baszamanx[0] = strconv.Itoa(yilint)
				baszamanx[1] = strconv.Itoa(ayint)
			}
			tut = baszamanx[0] + "-" + baszamanx[1] + "-" + baszamanx[2]
			sorgustr += " and tarih >= '" + tut + "'"
		} else {
			sorgustr += " and tarih >= '" + linkveri.BasTarih + "'"
		}
	}
	if linkveri.BitTarih != "" {
		sorgustr += " and tarih <= '" + linkveri.BitTarih + "'"
	}
	sorgustr += " order by id desc"

	var borclandirmalar []Borclandirma
	var odemeler []Odeme
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	rows, err := db.Query(sorgustr)
	hata(err)

	for rows.Next() {
		var bo Borclandirma
		err = rows.Scan(&bo.Id, &bo.Siteid, &bo.Blok, &bo.Daireno, &bo.Miktar, &bo.Kime, &bo.Zamangrup, &bo.Tarih, &bo.Gelirgiderid,
			&bo.Firmaid, &bo.Not, &bo.Zaman)
		hata(err)

		tutborc, err := strconv.ParseFloat(bo.Miktar, 64)
		hata(err)
		borctoplam += tutborc //Burada Borçlar Toplanıyor

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
			odemeler = append(odemeler, o)

			tutodeme, err := strconv.ParseFloat(o.Miktar, 64)
			hata(err)
			odemetoplam += tutodeme //Burada Ödemeler Toplanıyor
		}

		err = db.QueryRow("select ad from gelirgider where id = ?", bo.Gelirgiderid).Scan(&bo.Gelirgiderad)
		hata(err)
		if bo.Firmaid != -1 {
			err = db.QueryRow("select ad from firma where id = ?", bo.Firmaid).Scan(&bo.Firmaad)
			hata(err)
		}
		borclandirmalar = append(borclandirmalar, bo)
	}
	return borclandirmalar, odemeler, borctoplam, odemetoplam
}

func kisilereDagit(borclandirma []Borclandirma, kisiid int) []OzelRapor {
	var ozelrapor []OzelRapor
	for _, borc := range borclandirma {
		dairehareketleri := daireHareketlerGetirx(borc.Blok, borc.Daireno)
		if borc.Zamangrup == "1" { //anlık borçlar
			if kisiid == -1 { //Kişi Seçilmemişse
				ozelrapor = append(ozelrapor, anlikSorumluBulx(borc, dairehareketleri)) //Anlık Borç
			} else { //Kişiye göre ekleme yapacak
				or := anlikSorumluBulx(borc, dairehareketleri)
				if or.Kisi.Id == kisiid {
					ozelrapor = append(ozelrapor, or)
				}
			}
		} else if borc.Zamangrup == "2" { //aylık borçlar
			if kisiid == -1 { //Kişi Seçilmemişse
				ozelrapor = append(ozelrapor, aylikSorumluBulx(borc, dairehareketleri)...)
			} else { //Kişiye göre ekleme yapacak
				orlar := aylikSorumluBulx(borc, dairehareketleri)
				for _, or := range orlar {
					if or.Kisi.Id == kisiid {
						ozelrapor = append(ozelrapor, or)
					}
				}
			}
		}
	}
	return ozelrapor
}

func anlikSorumluBulx(borc Borclandirma, dairehareketleri []DaireHareket) OzelRapor {
	var onceki DaireHareket
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
	var ozelrapor OzelRapor
	kisiid := kimBulx(borc, onceki)
	ozelrapor.Borc = borc
	if kisiid != -1 {
		ozelrapor.Kisi = kisiGetir(kisiid)
	}
	return ozelrapor
}

func aylikSorumluBulx(borc Borclandirma, dairehareketleri []DaireHareket) []OzelRapor {
	var hareketler []DaireHareket
	var ozelrapor []OzelRapor

	//Burası tarihin ay kısmını bir arttırmaktadır.
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
		return ozelrapor
	} else if once && !orta {
		var or OzelRapor
		or.Borc = borc
		kisiid := kimBulx(borc, onceki)
		if kisiid != -1 {
			or.Kisi = kisiGetir(kisiid)
		}
		or.BorcGun = aycekimsayi
		miktar, err := strconv.ParseFloat(or.Borc.Miktar, 64)
		hata(err)
		or.BorcPayi = miktar
		ozelrapor = append(ozelrapor, or)
		return ozelrapor
	} else {
		if once {
			var or OzelRapor
			or.Borc = borc
			kisiid := kimBulx(borc, onceki)
			if kisiid != -1 {
				or.Kisi = kisiGetir(kisiid)
			}
			or.BorcGun = gunler[0]
			miktar, err := strconv.ParseFloat(or.Borc.Miktar, 64)
			hata(err)
			miktarx, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(or.BorcGun)/float64(aycekimsayi))*miktar), 64)
			hata(err)
			or.BorcPayi = miktarx
			if miktarx != 0 {
				ozelrapor = append(ozelrapor, or)
			}
		}
		for i, hareket := range hareketler {
			var or OzelRapor
			or.Borc = borc
			kisiid := kimBulx(borc, hareket)
			if kisiid != -1 {
				or.Kisi = kisiGetir(kisiid)
			}
			or.BorcGun = gunler[i+1]
			miktar, err := strconv.ParseFloat(or.Borc.Miktar, 64)
			hata(err)
			miktarx, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(or.BorcGun)/float64(aycekimsayi))*miktar), 64)
			hata(err)
			or.BorcPayi = miktarx
			ozelrapor = append(ozelrapor, or)
		}
	}

	var ormap = make(map[int]OzelRapor)
	for _, or := range ozelrapor {
		x := ormap[or.Kisi.Id]
		x.Borc = or.Borc
		x.Kisi = or.Kisi
		x.BorcPayi = or.BorcPayi
		x.BorcGun += or.BorcGun
		ormap[or.Kisi.Id] = x
	}
	var orapor []OzelRapor
	for _, rapor := range ormap {
		orapor = append(orapor, rapor)
	}
	return orapor
}

// -1 dönerse sorumlu yoktur demektir.
func kimBulx(borc Borclandirma, hareket DaireHareket) int {
	if borc.Kime == "1" { //Oturan
		if hareket.KiraciId == -1 {
			return hareket.MulksahibiId
		} else {
			return hareket.KiraciId
		}
	} else if borc.Kime == "2" { //Mülk Sahibi
		return hareket.MulksahibiId
	} else if borc.Kime == "3" { //Kiracı
		return hareket.KiraciId
	}
	return -1
}

func kisiGetir(id int) Kisi {
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	siteid, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	var kisi Kisi
	err = db.QueryRow("select * from kisi where siteid = ? and id = ?", siteid, id).
		Scan(&kisi.Id, &kisi.Siteid, &kisi.Ad, &kisi.Soyad, &kisi.Mail, &kisi.Telefon, &kisi.Cinsiyet, &kisi.Yas, &kisi.Zaman)
	hata(err)
	zamanint, err := strconv.ParseInt(kisi.Zaman, 10, 64)
	hata(err)
	sdizi := strings.Split(time.Unix(zamanint, 0).String(), " ")
	kisi.Zaman = sdizi[0] + " " + sdizi[1]
	return kisi
}

func daireHareketlerGetirx(blok, daireno string) []DaireHareket {
	var dairehareketler []DaireHareket
	db, err := sql.Open("sqlite3", "apartmanx.db")
	hata(err)
	defer db.Close()
	id, err := strconv.Atoi(aktifSiteGetir())
	hata(err)
	rows, err := db.Query("select * from dairehareket where siteid = ? and blok = ? and daireno = ? order by tarih asc", id, blok, daireno)
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
