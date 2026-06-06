package parser

import (
	"bytes"
	"context"
	"io"

	"github.com/ledongthuc/pdf"
)

type PDFExtractor struct{}

func (e *PDFExtractor) SupportedExtensions() []string {
	return []string{".pdf"}
}

func (e *PDFExtractor) Extract(ctx context.Context, reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	var text string
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		text += content + "\n"
	}

	return text, nil
}
