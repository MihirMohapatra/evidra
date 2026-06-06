package service

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

func NewXLSXGenerator() *excelize.File {
	return excelize.NewFile()
}

func XLSXBytes(f *excelize.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
