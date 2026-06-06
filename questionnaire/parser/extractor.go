package parser

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type Extractor interface {
	Extract(ctx context.Context, reader io.Reader) (string, error)
	SupportedExtensions() []string
}

func ExtractorFor(ext string) (Extractor, error) {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	switch ext {
	case "pdf":
		return &PDFExtractor{}, nil
	case "xlsx", "xls":
		return &XLSXExtractor{}, nil
	case "docx":
		return &DOCXExtractor{}, nil
	case "txt":
		return &TXTExtractor{}, nil
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}
