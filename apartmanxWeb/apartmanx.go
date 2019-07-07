package apartmanx

import (
	"net/http"
	"fmt"
	"os"
)

type Apartman struct {
	HeaderMenu    []Menu
	SiteEkleDurum bool
	Ad            string
	SiteAd        string
}

func logYaz(hata string) {
	hata += "\r\n"
	dosya, err := os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	dosya.Write([]byte(hata))
	dosya.Close()
}

func hata(err error) {
	if err != nil {
		logYaz(err.Error())
		panic(err)
	}
}

func apartman(w http.ResponseWriter, r *http.Request) {
	if oturumKontrol(w, r) {
		anasayfa(w, r)
	} else {
		giris(w, r)
	}
}

func Run() {
	openbrowser("http://localhost:4182")
	http.Handle("/root/", http.StripPrefix("/root/", http.FileServer(http.Dir("root"))))
	http.HandleFunc("/", apartman)
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/iletisim", iletisim)
	http.HandleFunc("/cikis", cikis)
	http.HandleFunc("/post", post)
	http.HandleFunc("/rapor", rapor)
	http.HandleFunc("/siteler/siteEkle", siteEkle)
	http.HandleFunc("/siteler/siteDuzenle", siteDuzenle)
	http.HandleFunc("/daireIslemleri/daireKayitDuzenleme", daireKayitDuzenleme)
	http.HandleFunc("/daireIslemleri/daireListesi", daireListesi)
	http.HandleFunc("/daireIslemleri/daireHareketleri", daireHareketleri)
	http.HandleFunc("/raporlar/tanimliRaporlar", tanimliRaporlar)
	http.HandleFunc("/raporlar/ozelRaporHazirla", ozelRaporHazirla)
	http.HandleFunc("/tanimlar/kisiTanimi", kisiTanimi)
	http.HandleFunc("/tanimlar/firmaTanimi", firmaTanimi)
	http.HandleFunc("/tanimlar/gelirGiderTanimi", gelirGiderTanimi)
	http.HandleFunc("/borclandirmaOdeme/borclandirma", borclandirma)
	http.HandleFunc("/borclandirmaOdeme/odeme", odeme)
	http.HandleFunc("/mail/mailGonderme", mailGonderme)
	err := http.ListenAndServe(":4182", nil)
	hata(err)

}
