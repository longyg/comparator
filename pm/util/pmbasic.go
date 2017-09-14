package pmutil

import "encoding/xml"

type PMBasic struct {
	XMLName xml.Name `xml:"PMBasic"`
	Xmlns2 string `xml:"ns2,attr"`
	Xmlns3 string `xml:"ns3,attr"`
	XmlnsXSI string `xml:xsi,attr`
	SchemaVersion string `xml:"schemaVersion,attr"`
	InterfaceVersion string`xml:"interfaceVersion, attr"`
	SchemaLocation string `xml:"schemaLocation, attr"`
	Adaptation Adaptation `xml:"Adaptation"`
	Product Product `xml:"Product"`
	PMClasses PMClasses `xml:"PMClasses"`
	Measurement []Measurement `xml:Measurement`
	AdminOfMeasurements *AdminOfMeasurements
}

type Adaptation struct {
	XMLName xml.Name `xml:"Adaptation"`
	Href string `xml:"href,attr"`
}

type Product struct {
	XMLName xml.Name `xml:"Product"`
	TypeName string `xml:"typeName,attr"`
	Vendor string `xml:"vendor,attr"`
	Release string `xml:release,attr`
	Description string
}

type PMClasses struct {
	XMLName xml.Name `xml:"PMClasses"`
	PMClassInfo []PMClassInfo `xml:"PMClassInfo"`
}

type PMClassInfo struct {
	XMLName xml.Name `xml:PMClassInfo`
	Description string
	Annotation []Annotation
	Name string `xml:"name,attr"`
	NameInOMeS string `xml:"nameInOMeS,attr"`
	Intransient string `xml:"intransient,attr"`
	Presentation string `xml:"presentation,attr"`
}

type Measurement struct {
	XMLName xml.Name `xml:"Measurement"`
	MeasuredTarget []MeasuredTarget
	Description string
	MeasuredIndicator []MeasuredIndicator
	SupportedMeasurementIntervals *SupportedMeasurementIntervals
	Annotation []Annotation
	MeasurementType string `xml:"measurementType,attr"`
	MeasurementTypeInOMeS string `xml:"measurementTypeInOMeS,attr"`
	Presentation string `xml:"presentation,attr"`
	DefaultInterval string `xml:"defaultInterval,attr"`
	NetworkElementMeasurementId string `xml:"networkElementMeasurementId,attr"`
	OmesFileGroupName string `xml:"omesFileGroupName,attr"`
	ResultsPerInterval string `xml:"resultsPerInterval,attr"`
	AomSupported string `xml:"aomSupported,attr"`
}

type MeasuredTarget struct {
	XMLName xml.Name `xml:"MeasuredTarget"`
	Hierarchy []Hierarchy
	Dimension string `xml:"dimension,attr"`
}

type Hierarchy struct {
	XMLName xml.Name `xml:"Hierarchy"`
	Classes string `xml:"classes,attr"`
}

type MeasuredIndicator struct {
	XMLName xml.Name `xml:"MeasuredIndicator"`
	Description string
	TimeAndObjectAggregationRule string
	TimeAndObjectAggregationFormula Formula
	ObjectAggregationRule string
	ObjectAggregationFormula Formula
	TimeAggregationRule string
	TimeAggregationFormula Formula
	Documentation MeasuredIndicatorDocumentation
	SupportedInProducts SupportedInProducts
	Annotation []Annotation

	Name string `xml:"name,attr"`
	NameInOMeS string `xml:"nameInOMeS,attr"`
	DataCollectionMethod string `xml:"dataCollectionMethod,attr"`
	Presentation string `xml:"presentation,attr"`
	Unit string `xml:"unit,attr"`
	DefaultValue string `xml:"defaultValue,attr"`
	MinValue string `xml:"minValue,attr"`
	MaxValue string `xml:"maxValue,attr"`
}

type MeasuredIndicatorDocumentation struct {
	Updated string
	Feature []Feature
	Standard []Standard
	NetworkElementVersion string `xml:"networkElementVersion,attr"`
	NetworkElementDevelopmentState string `xml:"networkElementDevelopmentState,attr"`
	OriginalRelease string `xml:"originalRelease,attr"`
	EndRelease string `xml:"endRelease,attr"`
}

type SupportedInProducts struct {
	ProductRef []ProductRef
}

type ProductRef struct {
	Annotation []Annotation
	TypeName string `xml:"typeName,attr"`
}

type Feature struct {
	Description string
	Id string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Status string `xml:"status,attr"`
}

type Standard struct {
	StandardName string
	StandardVersion string
}

type Formula struct {
	ValueMap FormulaValueMap
	StringFormula StringFormula
	Description string
}

type FormulaValueMap struct {
	Variable string `xml:"variable,attr"`
	Value string `xml:"value,attr"`
	ReplacementValue string `xml:"replacementValue,attr"`
}

type StringFormula struct {
	Language string `xml:"language,attr"`
}

type SupportedMeasurementIntervals struct {
	XMLName xml.Name `xml:"SupportedMeasurementIntervals"`
	SupportedMeasurementInterval []string
	Documentation string
	RecommendedInterval string `xml:"recommendedInterval,attr"`
}

type Annotation struct {
	XMLName xml.Name `xml:"Annotation"`
	Elements []AnnotationElement `xml:"elements"`
	Type Reference `xml:"type"`
}

type Reference struct {
	Href string `xml:"href,attr"`
}

type AnnotationElement struct {
	Name string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type AdminOfMeasurements struct {
	XMLName xml.Name `xml:"AdminOfMeasurements"`
	SupportedMeasurementIntervals *SupportedMeasurementIntervals
}