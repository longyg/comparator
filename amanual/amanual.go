package amanual

import (
	"encoding/xml"
)

type Manual struct {
	XMLName xml.Name `xml:"AlarmDescription"`
	XMIVersion string `xml:"version,attr"`
	XMIXmlns string `xml:"xmi,attr"`
	Xmlns string `xml:"com.nokia.oss.fm.fmbasic,attr"`
	PatchLevel string `xml:"patchLevel,attr"`
	SpecificProblem string `xml:"specificProblem,attr"`
	AlarmText string `xml:"alarmText,attr"`
	ProbableCause string `xml:"probableCause,attr"`
	SupplementaryInformationFields string `xml:"supplementaryInformationFields,attr"`
	PerceivedSeverityInfo string `xml:"perceivedSeverityInfo,attr"`
	AlarmType string `xml:"alarmType,attr"`
	Meaning string `xml:"meaning,attr"`
	Instructions string `xml:"instructions,attr"`
	Cancelling string `xml:"cancelling,attr"`
	Adaptation Adaptation `xml:"Adaptation"`
}

type Adaptation struct {
	XMLName xml.Name `xml:"Adaptation"`
	Href string `xml:"href,attr"`
}

