package utils

import (
	"log"
	"runtime"
	"strings"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

/**
 * example using
 * utils.ExportToPdf("https://google.com")
 */
func ExportToPdf(url string) []byte {
	if runtime.GOOS == "windows" {
		wkhtmltopdf.SetPath("C:/Program Files/wkhtmltopdf/bin/wkhtmltopdf.exe")
	}
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println("Error creating PDF generator:", err)
		return nil
	}

	// Set global options
	pdfg.Dpi.Set(350)

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPage(url)
	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Println("Error creating PDF document:", err)
		return nil
	}

	return pdfg.Bytes()
}

func ExportHtmlToPdf(html string, basePath string) []byte {
	if runtime.GOOS == "windows" {
		wkhtmltopdf.SetPath("C:/Program Files/wkhtmltopdf/bin/wkhtmltopdf.exe")
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println("Error creating PDF generator:", err)
		return nil
	}

	pdfg.Dpi.Set(350)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	page := wkhtmltopdf.NewPageReader(strings.NewReader(html))
	page.DisableJavascript.Set(true)
	page.LoadMediaErrorHandling.Set("ignore")
	page.LoadErrorHandling.Set("ignore")
	if basePath != "" {
		page.Allow.Set(basePath)
	}
	
	pdfg.AddPage(page)

	// Log info
	log.Printf("Creating PDF from HTML (%d bytes). BasePath: %s", len(html), basePath)

	err = pdfg.Create()
	if err != nil {
		log.Println("Error creating PDF document from HTML:", err)
		return nil
	}

	return pdfg.Bytes()
}
