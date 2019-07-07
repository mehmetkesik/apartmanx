function daireDegisim(secimid) {
    var blokdaireno = document.getElementById("blokdaireno");
    var blok = document.getElementById("blok");
    var daireno = document.getElementById("daireno");
    if (secimid == 1) {
        blokdaireno.style.display = "flex";
        blok.removeAttribute("disabled");
        daireno.removeAttribute("disabled");
    } else if (secimid == 2) {
        blokdaireno.style.display = "none";
        blok.setAttribute("disabled", "disabled");
        daireno.setAttribute("disabled", "disabled");
    }
}

function tarihDegisim(secimid) {
    var anlik = document.getElementById("anlik");
    var aylik = document.getElementById("aylik");
    var tarih1 = document.getElementById("tarih1");
    var tarih2 = document.getElementById("tarih2");
    if (secimid == 1) {
        anlik.style.display = "flex";
        aylik.style.display = "none";
        tarih1.removeAttribute("disabled");
        tarih2.setAttribute("disabled", "disabled");
    } else if (secimid == 2) {
        anlik.style.display = "none";
        aylik.style.display = "flex";
        tarih1.setAttribute("disabled", "disabled");
        tarih2.removeAttribute("disabled");
    }
}

function raporal() {
    window.open("http://" + window.location.host + "/rapor" + window.location.search, "_blank");
}