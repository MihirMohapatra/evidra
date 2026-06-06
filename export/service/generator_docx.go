package service

import (
	"archive/zip"
	"bytes"
	"fmt"
)

type DOCXGenerator struct {
	title string
	paras []string
}

func NewDOCXGenerator() *DOCXGenerator {
	return &DOCXGenerator{}
}

func (d *DOCXGenerator) AddTitle(text string) {
	d.title = text
}

func (d *DOCXGenerator) AddParagraph(text string) {
	d.paras = append(d.paras, text)
}

func (d *DOCXGenerator) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>`

	if d.title != "" {
		content += fmt.Sprintf(`
    <w:p>
      <w:pPr><w:pStyle w:val="Title"/></w:pPr>
      <w:r><w:t>%s</w:t></w:r>
    </w:p>`, xmlEscape(d.title))
	}

	for _, p := range d.paras {
		content += fmt.Sprintf(`
    <w:p>
      <w:r><w:t>%s</w:t></w:r>
    </w:p>`, xmlEscape(p))
	}

	content += `
  </w:body>
</w:document>`

	if err := addZipFile(w, "word/document.xml", content); err != nil {
		return nil, err
	}

	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
	if err := addZipFile(w, "_rels/.rels", rels); err != nil {
		return nil, err
	}

	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
	if err := addZipFile(w, "[Content_Types].xml", contentTypes); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func addZipFile(w *zip.Writer, name, content string) error {
	f, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	return err
}

func xmlEscape(s string) string {
	var buf bytes.Buffer
	for _, c := range s {
		switch c {
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '"':
			buf.WriteString("&quot;")
		case '\'':
			buf.WriteString("&apos;")
		default:
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
