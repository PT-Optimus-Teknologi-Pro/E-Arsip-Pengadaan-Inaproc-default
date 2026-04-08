package utils

import (
	"log"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

const path = "/usr/local/bin/wkhtmltopdf"
/**
 * example using
 * utils.ExportToPdf("https://google.com", "./simplesample.pdf")
 */
func ExportToPdf(url string) []byte {
	// Create new PDF generator
  pdfg, err := wkhtmltopdf.NewPDFGenerator()
  if err != nil {
    log.Fatal(err)
  }

  // Set global options
  pdfg.Dpi.Set(350)
  // pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
  // pdfg.Grayscale.Set(true)
  // pdfg.MarginBottom.Set(0)
  // pdfg.MarginTop.Set(0)
  // pdfg.MarginLeft.Set(0)
  // pdfg.MarginRight.Set(0)

  // Create a new input page from an URL
  page := wkhtmltopdf.NewPage(url)
  // Add to document
  pdfg.AddPage(page)

  // Create PDF document in internal buffer
  err = pdfg.Create()
  if err != nil {
    log.Fatal(err)
  }

  return pdfg.Bytes()
  // Write buffer contents to file on disk
  // err = pdfg.WriteFile(outputPath)
  // if err != nil {
  //   log.Fatal(err)
  // }
}
