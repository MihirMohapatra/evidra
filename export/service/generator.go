package service

import (
	"bytes"

	gofpdf "github.com/jung-kurt/gofpdf"
)

func NewPDFGenerator() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	return pdf
}

func PDFBytes(pdf *gofpdf.Fpdf) ([]byte, error) {
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
