package parser

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"strings"
)

type DOCXExtractor struct{}

func (e *DOCXExtractor) SupportedExtensions() []string {
	return []string{".docx"}
}

func (e *DOCXExtractor) Extract(ctx context.Context, reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	var documentXML []byte
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			documentXML, err = io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			break
		}
	}

	if documentXML == nil {
		return "", nil
	}

	return extractTextFromDocxXML(documentXML), nil
}

type wDoc struct {
	Body wBody `xml:"body"`
}

type wBody struct {
	Paragraphs []wParagraph `xml:"p"`
}

type wParagraph struct {
	Runs []wRun `xml:"r"`
}

type wRun struct {
	Text string `xml:"t"`
}

func extractTextFromDocxXML(data []byte) string {
	var doc wDoc
	if err := xml.Unmarshal(data, &doc); err != nil {
		return ""
	}

	var texts []string
	for _, p := range doc.Body.Paragraphs {
		var line string
		for _, r := range p.Runs {
			line += r.Text
		}
		texts = append(texts, strings.TrimSpace(line))
	}

	return strings.Join(texts, "\n")
}
