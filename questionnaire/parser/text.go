package parser

import (
	"context"
	"io"
)

type TXTExtractor struct{}

func (e *TXTExtractor) SupportedExtensions() []string {
	return []string{".txt"}
}

func (e *TXTExtractor) Extract(ctx context.Context, reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
