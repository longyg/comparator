package pmutil

import "encoding/xml"

type CompareResult struct {
	XMLName xml.Name `xml:"CompareResult"`
	Deviation []Deviation
}

type Deviation struct {
	XMLName xml.Name `xml:"Deviation"`
	Level string `xml:"level,attr"`
	PmClass string `xml:"pmClass,attr"`
	MeasurementType string `xml:"measurementType,attr"`
	Counter string `xml:"counter,attr"`
	Type string `xml:"deviationCategory,attr"`
	Message string `xml:"message"`
}
