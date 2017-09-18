package pmutil

import "encoding/xml"

type CompareResult struct {
	XMLName xml.Name `xml:"CompareResult"`
	Deviation []Deviation
}

type Deviation struct {
	XMLName xml.Name `xml:"Deviation"`
	Level string `xml:"level,attr"`
	DeviationCategory string `xml:"deviationCategory,attr"`
	Message string `xml:"message"`
	PmClass *string `xml:"pmClass,attr"`
	MeasurementType *string `xml:"measurementType,attr"`
	Dimension *string `xml:"dimension,attr"`
	Counter *string `xml:"counter,attr"`
}
