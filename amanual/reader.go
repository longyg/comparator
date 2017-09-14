package amanual

import (
	"io"
	"encoding/xml"
	"io/ioutil"
	"strings"
	"os"
	"fmt"
)

func readMan(reader io.Reader) (Manual, error) {
	var man Manual
	err := xml.NewDecoder(reader).Decode(&man)
	if err != nil {
		return man, err
	}
	return man, nil
}

func ReadMan(filePath string) (Manual, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Open man file:", filePath, "failed")
		return Manual{}, err
	}
	return readMan(file)
}

func ReadMans(manDir string) map[string]Manual {
	files, _ := ioutil.ReadDir(manDir)

	manuals := make(map[string]Manual)

	for _, f := range files {
		filepath := strings.Join([]string{manDir, "\\", f.Name()}, "")
		man, err := ReadMan(filepath)

		if err != nil {
			fmt.Println("Read man file", f.Name(), "error:", err)
			continue
		}
		manuals[f.Name()] = man
	}
	return manuals
}
