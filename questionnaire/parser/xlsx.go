package parser

import (
	"bytes"
	"context"
	"io"

	"github.com/xuri/excelize/v2"
)

type XLSXExtractor struct{}

func (e *XLSXExtractor) SupportedExtensions() []string {
	return []string{".xlsx", ".xls"}
}

func (e *XLSXExtractor) Extract(ctx context.Context, reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	var text string
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		for _, row := range rows {
			for i, cell := range row {
				if i > 0 {
					text += "\t"
				}
				text += cell
			}
			text += "\n"
		}
	}

	return text, nil
}
