package pmutil

import (
	"os"
	"fmt"
	"io"
	"encoding/xml"
)

func ReadPmb(filePath string) (PMBasic, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Open pmb file:", filePath, "failed")
		return PMBasic{}, err
	}
	return readPmb(file)
}

func readPmb(reader io.Reader) (PMBasic, error) {
	var pmb PMBasic
	err := xml.NewDecoder(reader).Decode(&pmb)
	if err != nil {
		return pmb, err
	}
	return pmb, nil
}
