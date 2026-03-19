package packager

import (
	"archive/zip"
	"businesscard-wallet-backend/ports"
	"bytes"
	"fmt"
)

type ZipPackager struct{}

func NewZipPackager() *ZipPackager {
	return &ZipPackager{}
}

func (p *ZipPackager) Package(files []ports.PackageFile) ([]byte, error) {
	var buffer bytes.Buffer

	zipWriter := zip.NewWriter(&buffer)

	for _, file := range files {
		entryWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			_ = zipWriter.Close()
			return nil, fmt.Errorf("create zip entry %s: %w", file.Name, err)
		}

		if _, err := entryWriter.Write(file.Data); err != nil {
			_ = zipWriter.Close()
			return nil, fmt.Errorf("write zip entry %s: %w", file.Name, err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("close zip writer: %w", err)
	}

	return buffer.Bytes(), nil
}
