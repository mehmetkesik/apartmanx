package apartmanx

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func html2pdf(sourcelink, targerpdf string) {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	hata(err)

	// Add one page from an URL
	pdfg.AddPage(wkhtmltopdf.NewPage(sourcelink))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	hata(err)

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(targerpdf)
	hata(err)

	// Output: Done
}

/*
direk html kodlarını pdf'e çevirme
html := "<html>Hi</html>"
pdfgen.AddPage(NewPageReader(strings.NewReader(html)))
*/
