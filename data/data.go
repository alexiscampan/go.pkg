// Package data defines methods to transform data
package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// AsCSVBuf to get string data as csv format
func AsCSVBuf(entry []string) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	err := writer.Write(entry)
	if err != nil {
		log.Fatal("Entry data incompatible type", zap.Error(err))
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatal("Cannot write to buffer", zap.Error(err))
	}
	return buffer, nil
}

// ToCSVFile to write buffer to csv file
func ToCSVFile(buf bytes.Buffer, filename string) error {
	f, err := os.Create(filepath.Clean(filename))
	defer func() {
		if err2 := f.Close(); err2 != nil {
			if err == nil {
				err = fmt.Errorf("failed to close file")
			}
		}
	}()
	if err != nil {
		log.Fatal("Cannot create file ", zap.Error(err))
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
