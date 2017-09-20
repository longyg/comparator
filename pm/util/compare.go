package pmutil

import (
	"fmt"
	"os"
	"strings"
	"encoding/xml"
	"reflect"
)

var logFile *os.File
var resultFile *os.File
var result CompareResult
var deviations []Deviation

func Compare(beforePmb string, afterPmb string, logFilePath string, resultPath string) {
	beforePmbasic, err := ReadPmb(beforePmb)
	if nil != err {
		fmt.Println("error while reading pmb: " + beforePmb, err)
		return
	}
	afterPmbasic, err := ReadPmb(afterPmb)
	if nil != err {
		fmt.Println("error while reading pmb: " + afterPmb, err)
		return
	}

	logFile, _ = os.Create(logFilePath)
	resultFile, _ = os.Create(resultPath)
	defer logFile.Close()
	defer resultFile.Close()

	result = CompareResult{}

	compare(beforePmbasic, afterPmbasic)

	result.Deviation = deviations

	writeResult()
}

func compare(beforePmb PMBasic, afterPmb PMBasic) {
	compareBasicInfo(beforePmb, afterPmb)
	comparePMClasses(beforePmb, afterPmb)
	compareMeasurements(beforePmb, afterPmb)
	compareAdminOfMeasurements(beforePmb.AdminOfMeasurements, afterPmb.AdminOfMeasurements)
}

type Context struct {
	PMClass string
	MeasurementType string
	Counter string
	Dimension string
}

func compareBasicInfo(beforePmb PMBasic, afterPmb PMBasic) {
	pmBasicFields := []interface{}{
		FieldDef{"Xmlns3", "", "xmlns:ns3", false},
		FieldDef{"XmlnsXSI", "", "xmlns:xsi", false},
		FieldDef{"SchemaVersion", "", "schemaVersion", false},
		FieldDef{"InterfaceVersion", "", "interfaceVersion", false},
		FieldDef{"SchemaLocation", "", "schemaLocation", false},
	}
	compareFields(beforePmb, afterPmb, pmBasicFields)

	adaptationFields := []interface{}{
		FieldDef{"Href", "", "href", false},
	}
	compareFields(beforePmb.Adaptation, afterPmb.Adaptation, adaptationFields)

	productFields := []interface{}{
		FieldDef{"TypeName", "", "typeName", false},
		FieldDef{"Vendor", "", "vendor", false},
		FieldDef{"Release", "", "release", false},
		FieldDef{"Description", "", "Description", false},
	}
	compareFields(beforePmb.Product, afterPmb.Product, productFields)
}

func compareFields(before interface{}, after interface{}, fields interface{}) {
	bv := reflect.ValueOf(before)
	elementType := getElementTypeName(bv)

	av := reflect.ValueOf(after)
	category := fmt.Sprint(elementType, " Deviation")

	ft := reflect.TypeOf(fields)
	if ft.Kind() != reflect.Slice {
		panic("'fields' argument is not slice")
	}

	fv := reflect.ValueOf(fields)
	for i :=0; i < fv.Len(); i++ {
		fve := fv.Index(i)
		fvet := fve.Elem().Type().Name()
		if fvet == "FieldDef" {
			fieldName := getFieldStringValue(fve.Elem(), "FieldName")
			defaultValue := getFieldStringValue(fve.Elem(), "DefaultValue")
			xmlFieldName := getFieldStringValue(fve.Elem(), "XmlFieldName")
			handleEmpty := getFieldBoolValue(fve.Elem(), "HandleEmpty")

			bFieldValue := getFieldStringValue(bv, fieldName)
			aFieldValue := getFieldStringValue(av, fieldName)

			if handleEmpty {
				if bFieldValue == "" && defaultValue != "" {
					bFieldValue = defaultValue
				}
				if aFieldValue == "" && defaultValue != "" {
					aFieldValue = defaultValue
				}
			}

			fieldInfo := fmt.Sprint(elementType, " [", xmlFieldName, "]")
			compareField(fieldInfo, bFieldValue, aFieldValue, true, category, Context{})
		} else {
			fmt.Println("Error: Invalid field", fve)
			panic("Invalid field")
		}
	}
}

func compareNumber(element string, before interface{}, after interface{}, context Context)  {
	bt := reflect.TypeOf(before)
	at := reflect.TypeOf(after)

	if bt.Kind() != reflect.Slice {
		panic("before argument is not slice")
	}
	if at.Kind() != reflect.Slice {
		panic("after argument is not slice")
	}
	lenBefore := reflect.ValueOf(before).Len()
	lenAfter := reflect.ValueOf(after).Len()
	if lenBefore != lenAfter {
		message := fmt.Sprint("Number of ", element, " is not equal, before: {", lenBefore, "}, after: {", lenAfter, "}")
		category := fmt.Sprint(element, " Number Deviation")
		logDeviation(category, message, context)
	}
}

func compareElements(before interface{}, after interface{}, keyFieldName string, fields interface{}, context Context) {
	bt := reflect.TypeOf(before)
	at := reflect.TypeOf(after)

	if bt.Kind() != reflect.Slice {
		panic("before argument is not slice")
	}
	if at.Kind() != reflect.Slice {
		panic("after argument is not slice")
	}

	bv := reflect.ValueOf(before)
	av := reflect.ValueOf(after)

	compareBeforeWithAfter(bv, av, keyFieldName, fields, context)
	compareAfterWithBefore(av, bv, keyFieldName, context)
}

func compareBeforeWithAfter(bv reflect.Value, av reflect.Value, keyFieldName string, fields interface{}, context Context) {
	for i := 0; i < bv.Len(); i++ {
		be := bv.Index(i)
		bTypeName := be.Type().Name()
		bFieldValue := getFieldValue(be, keyFieldName)

		isFound := false

		for j := 0; j < av.Len(); j++ {
			ae := av.Index(j)
			aFieldValue := getFieldValue(ae, keyFieldName)

			if bFieldValue == aFieldValue {
				isFound = true
				ctx := setContext(context, bTypeName, bFieldValue)
				compareElementFields(be, ae, keyFieldName, fields, ctx)
				break
			}
		}
		if !isFound {
			message := fmt.Sprint(fmt.Sprint(bTypeName, " [", bFieldValue, "] is missing in exported PMB"))
			category := fmt.Sprint("Missing ", bTypeName)
			logDeviation(category, message, context)
		}
	}
}

func setContext(context Context, typeName string, fieldValue string) Context {
	if typeName == "PMClassInfo" {
		context.PMClass = fieldValue
	}
	if typeName == "Measurement" {
		context.MeasurementType = fieldValue
	}
	if typeName == "MeasuredTarget" {
		context.Dimension = fieldValue
	}
	if typeName == "MeasuredIndicator" {
		context.Counter = fieldValue
	}
	return context
}

func compareAfterWithBefore(av reflect.Value, bv reflect.Value, keyFieldName string, context Context) {
	for i := 0; i < av.Len(); i++ {
		ae := av.Index(i)
		aTypeName := ae.Type().Name()
		aFieldValue := getFieldValue(ae, keyFieldName)

		isFound := false

		for j := 0; j < bv.Len(); j++ {
			be := bv.Index(j)
			bFieldValue := getFieldValue(be, keyFieldName)

			if aFieldValue == bFieldValue {
				isFound = true
				break
			}
		}
		if !isFound {
			message := fmt.Sprint(fmt.Sprint(aTypeName, " [", aFieldValue, "] is adding in exported PMB"))
			category := fmt.Sprint("Adding ", aTypeName)
			logDeviation(category, message, context)
		}
	}
}

func getFieldValue(v reflect.Value, fieldName string) string {
	fieldValue := ""
	if v.Kind() == reflect.String {
		fieldValue = v.String()
	} else {
		fieldValue = v.FieldByName(fieldName).String()
	}
	return fieldValue
}

type FieldDef struct {
	FieldName string
	DefaultValue string
	XmlFieldName string
	HandleEmpty bool
}

func getFieldStringValue(v reflect.Value, fieldName string) string {
	return v.FieldByName(fieldName).String()
}

func getFieldBoolValue(v reflect.Value, fieldName string) bool {
	return v.FieldByName(fieldName).Bool()
}

func getElementTypeName(v reflect.Value) string {
	return v.Type().Name()
}

func compareElementFields(bv reflect.Value, av reflect.Value, keyFieldName string, fields interface{}, context Context) {
	elementType := getElementTypeName(bv)
	bIdentiferFieldValue := getFieldStringValue(bv, keyFieldName)

	ft := reflect.TypeOf(fields)
	if ft.Kind() != reflect.Slice {
		panic("'fields' argument is not slice")
	}
	fv := reflect.ValueOf(fields)
	for i :=0; i < fv.Len(); i++ {
		fve := fv.Index(i)
		fvet := fve.Elem().Type().Name()
		if fvet == "FieldDef" {
			fieldName := getFieldStringValue(fve.Elem(), "FieldName")
			defaultValue := getFieldStringValue(fve.Elem(), "DefaultValue")
			xmlFieldName := getFieldStringValue(fve.Elem(), "XmlFieldName")
			handleEmpty := getFieldBoolValue(fve.Elem(), "HandleEmpty")

			bFieldValue := getFieldStringValue(bv, fieldName)
			aFieldValue := getFieldStringValue(av, fieldName)
			if elementType == "Measurement" && fieldName == "DefaultInterval" && bIdentiferFieldValue == "UB6FR"  {
				fmt.Println("==============>", bFieldValue)
				fmt.Println("==============>", aFieldValue)
			}

			if handleEmpty {
				if "" == bFieldValue {
					if "" == defaultValue {
						bFieldValue = bIdentiferFieldValue
					} else {
						bFieldValue = defaultValue
					}
				}
				if "" == aFieldValue {
					if "" == defaultValue {
						aFieldValue = getFieldStringValue(av, keyFieldName)
					} else {
						aFieldValue = defaultValue
					}
				}
			}

			compareFieldValue(elementType, bIdentiferFieldValue, bFieldValue, aFieldValue, xmlFieldName, context)

			if elementType == "Measurement" && xmlFieldName == "defaultInterval" && bIdentiferFieldValue == "UB6FR"  {
				fmt.Println("UB6FR")
				compareFieldValue(elementType, bIdentiferFieldValue, bFieldValue, aFieldValue, xmlFieldName, context)
			}
		} else {
			fmt.Println("Error: Invalid field", fve)
			panic("Invalid field")
		}
	}
}

func compareFieldValue(elementType string, identifer string, bFieldValue string, aFieldValue string, xmlFieldName string, context Context) {
	if elementType == "Measurement" && xmlFieldName == "defaultInterval" && identifer == "UB6FR"  {
		fmt.Println("==============>", bFieldValue)
		fmt.Println("==============>", aFieldValue)
	}
	fieldInfo := fmt.Sprint(elementType, " [", identifer, "] ", xmlFieldName)
	category := fmt.Sprint(elementType, " Deviation")
	if elementType == "MeasuredIndicator" {
		category = "Counter Deviation"
	}
	compareField(fieldInfo, bFieldValue, aFieldValue, true, category, context)
}

func comparePMClasses(beforePmb PMBasic, afterPmb PMBasic)  {
	ctx := Context{}
	compareNumber("PMClassInfo", beforePmb.PMClasses.PMClassInfo, afterPmb.PMClasses.PMClassInfo, ctx)

	fields := []interface{} {
		FieldDef{"NameInOMeS", "", "nameInOMeS", true},
		FieldDef{"Intransient", "", "intransient", false},
		FieldDef{"Presentation", "", "presentation", true},
		FieldDef{"Description", "", "Description", false},
	}
	compareElements(beforePmb.PMClasses.PMClassInfo, afterPmb.PMClasses.PMClassInfo, "Name", fields, ctx)
}

func compareAdminOfMeasurements(before *AdminOfMeasurements, after *AdminOfMeasurements)  {
	category := "AoM Config Deviation"
	context := Context{}
	if nil != before && nil != after {
		compareSupportedMeasurementIntervals(before.SupportedMeasurementIntervals, after.SupportedMeasurementIntervals, Context{})
	} else if nil != before && nil == after {
		message := fmt.Sprint("AdminOfMeasurements is missing in exported PMB")
		logDeviation(category, message, context)
	} else if nil == before && nil != after {
		message := fmt.Sprint("AdminOfMeasurements is adding in exported PMB")
		logDeviation(category, message, context)
	}
}

func compareMeasurements(beforePmb PMBasic, afterPmb PMBasic) {
	ctx := Context{}
	compareNumber("Measurement", beforePmb.Measurement, afterPmb.Measurement, ctx)

	measurementFields := []interface{} {
		FieldDef{"MeasurementTypeInOMeS", "", "measurementTypeInOMeS", true},
		FieldDef{"Presentation", "", "presentation", true},
		FieldDef{"DefaultInterval", "", "defaultInterval", true},
		FieldDef{"NetworkElementMeasurementId", "", "networkElementMeasurementId", true},
		FieldDef{"OmesFileGroupName", "", "omesFileGroupName", true},
		FieldDef{"ResultsPerInterval", "", "resultsPerInterval", true},
		FieldDef{"AomSupported", "true", "aomSupported", true},
		FieldDef{"Description", "", "Description", true},
	}
	compareElements(beforePmb.Measurement, afterPmb.Measurement, "MeasurementType", measurementFields, ctx);

	compareMeasurementChildElements(beforePmb, afterPmb)
}

func compareMeasurementChildElements(beforePmb PMBasic, afterPmb PMBasic) {
	for _, beforeMeas := range beforePmb.Measurement {
		for _, afterMeas := range afterPmb.Measurement {
			if beforeMeas.MeasurementType == afterMeas.MeasurementType {

				measCtx := Context {"", beforeMeas.MeasurementType, "", ""}
				compareNumber("MeasuredTarget", beforeMeas.MeasuredTarget, afterMeas.MeasuredTarget, measCtx)
				compareElements(beforeMeas.MeasuredTarget, afterMeas.MeasuredTarget, "Dimension", []interface{}{}, measCtx)

				for _, bMt := range beforeMeas.MeasuredTarget {
					for _, aMt := range afterMeas.MeasuredTarget {
						if bMt.Dimension == aMt.Dimension {
							mTCtx := Context{"", beforeMeas.MeasurementType, "", bMt.Dimension}
							compareNumber("Hierarchy", bMt.Hierarchy, aMt.Hierarchy, mTCtx)
							compareElements(bMt.Hierarchy, aMt.Hierarchy, "Classes", []interface{}{}, mTCtx)
						}
					}
				}

				compareNumber("Counter", beforeMeas.MeasuredIndicator, afterMeas.MeasuredIndicator, measCtx)
				counterFields := []interface{} {
					FieldDef{"NameInOMeS", "", "nameInOMeS", true},
					FieldDef{"Description", "", "Description", false},
					FieldDef{"DataCollectionMethod", "", "dataCollectionMethod", false},
					FieldDef{"Presentation", "", "presentation", true},
					FieldDef{"Unit", "", "unit", false},
					FieldDef{"DefaultValue", "", "defaultValue", false},
					FieldDef{"MinValue", "", "minValue", false},
					FieldDef{"MaxValue", "", "maxValue", false},
				}
				compareElements(beforeMeas.MeasuredIndicator, afterMeas.MeasuredIndicator, "Name", counterFields, measCtx)

				for _, bCounter := range beforeMeas.MeasuredIndicator {
					for _, aCounter := range afterMeas.MeasuredIndicator {
						if bCounter.Name == aCounter.Name {
							context := Context{"", beforeMeas.MeasurementType, bCounter.Name, ""}
							compareAggRules(bCounter, aCounter, context)
							compareFormula(bCounter.TimeAndObjectAggregationFormula, aCounter.TimeAndObjectAggregationFormula)
							//compareDocumentation(bCounter.Documentation, aCounter.Documentation, context)
							compareSupportedInProducts(bCounter.SupportedInProducts, aCounter.SupportedInProducts, context)
						}
					}
				}

				compareSupportedMeasurementIntervals(beforeMeas.SupportedMeasurementIntervals, afterMeas.SupportedMeasurementIntervals, measCtx)
			}
		}
	}
}

func compareAggRules(before MeasuredIndicator, after MeasuredIndicator, context Context)  {
	category := "Agg Rule Deviation"
	if "" != before.TimeAndObjectAggregationRule && "" != after.TimeAndObjectAggregationRule {
		compareField("TimeAndObjectAggregationRule", before.TimeAndObjectAggregationRule, after.TimeAndObjectAggregationRule, true, category, context)
	} else if "" != before.TimeAndObjectAggregationRule && "" == after.TimeAndObjectAggregationRule {
		if "" != after.TimeAggregationRule && "" != after.ObjectAggregationRule {
			compareField("TimeAggregationRule", before.TimeAndObjectAggregationRule, after.TimeAggregationRule, true, category, context)
			compareField("ObjectAggregationRule", before.TimeAndObjectAggregationRule, after.ObjectAggregationRule, true, category, context)
		} else {
			if "" == after.TimeAggregationRule {
				message := fmt.Sprint("TimeAggregationRule is missing in exported PMB")
				category = "TimeAggregationRule Not Exported"
				logDeviation(category, message, context)
			}
			if "" == after.ObjectAggregationRule {
				message := fmt.Sprint("ObjectAggregationRule is missing in exported PMB")
				category = "ObjectAggregationRule Not Exported"
				logDeviation(category, message, context)
			}
		}
	} else if "" == before.TimeAndObjectAggregationRule && "" != after.TimeAndObjectAggregationRule {
		if "" != before.TimeAggregationRule && "" != before.ObjectAggregationRule {
			compareField("TimeAggregationRule", before.TimeAggregationRule, after.TimeAndObjectAggregationRule, true, category, context)
			compareField("ObjectAggregationRule", before.ObjectAggregationRule, after.TimeAndObjectAggregationRule, true, category, context)
		} else {

			if "" == before.TimeAggregationRule {
				message := fmt.Sprint("TimeAggregationRule is not defined in import PMB")
				category = "TimeAggregationRule Not Defined"
				logDeviation(category, message, context)
			}
			if "" == before.ObjectAggregationRule {
				message := fmt.Sprint("ObjectAggregationRule is not defined in import PMB")
				category = "ObjectAggregationRule Not Defined"
				logDeviation(category, message, context)
			}
		}
	} else {
		compareField("TimeAggregationRule", before.TimeAggregationRule, after.TimeAggregationRule, true, category, context)
		compareField("ObjectAggregationRule", before.ObjectAggregationRule, after.ObjectAggregationRule, true, category, context)
	}
}

func compareSupportedInProducts(before *SupportedInProducts, after *SupportedInProducts, context Context)  {
	category := "SupportedInProducts Deviation"
	if nil != before && nil != after {
		//
	} else if nil != before && nil == after {
		message := fmt.Sprint("SupportedInProducts is missing in exported PMB")
		logDeviation(category, message, context)
	} else if nil == before && nil != after {
		message := fmt.Sprint("SupportedInProducts is adding in exported PMB")
		logDeviation(category, message, context)
	}
}

func compareDocumentation(before *MeasuredIndicatorDocumentation, after *MeasuredIndicatorDocumentation, context Context)  {
	category := "Counter Documentation Deviation"
	if nil != before && nil != after {
		//
	} else if nil != before && nil == after {
		message := fmt.Sprint("Counter Documentation is missing in exported PMB")
		logDeviation(category, message, context)
	} else if nil == before && nil != after {
		message := fmt.Sprint("Counter Documentation is adding in exported PMB")
		logDeviation(category, message, context)
	}
}

func compareFormula(before Formula, after Formula)  {

}

func compareSimpleElements(elementName string, before interface{}, after interface{}, context Context) {
	bt := reflect.TypeOf(before)
	at := reflect.TypeOf(after)

	if bt.Kind() != reflect.Slice {
		panic("before argument is not slice")
	}
	if at.Kind() != reflect.Slice {
		panic("after argument is not slice")
	}

	bv := reflect.ValueOf(before)
	av := reflect.ValueOf(after)

	for i := 0; i < bv.Len(); i++ {
		be := bv.Index(i)
		bFieldValue := be.String()

		isFound := false

		for j := 0; j < av.Len(); j++ {
			ae := av.Index(j)
			aFieldValue := ae.String()

			if bFieldValue == aFieldValue {
				isFound = true
				break
			}
		}
		if !isFound {
			message := fmt.Sprint(fmt.Sprint(elementName, " [", bFieldValue, "] is missing in exported PMB"))
			category := fmt.Sprint("Missing ", elementName)
			logDeviation(category, message, context)
		}
	}

	for i := 0; i < av.Len(); i++ {
		ae := av.Index(i)
		aFieldValue := ae.String()

		isFound := false

		for j := 0; j < bv.Len(); j++ {
			be := bv.Index(j)
			bFieldValue := be.String()

			if aFieldValue == bFieldValue {
				isFound = true
				break
			}
		}
		if !isFound {
			message := fmt.Sprint(fmt.Sprint(elementName, " [", aFieldValue, "] is adding in exported PMB"))
			category := fmt.Sprint("Adding ", elementName)
			logDeviation(category, message, context)
		}
	}
}

func compareSupportedMeasurementIntervals(before *SupportedMeasurementIntervals, after *SupportedMeasurementIntervals, context Context)  {
	category := "AoM Config Deviation"
	if nil != before && nil != after {
		compareField("Documentation", before.Documentation, after.Documentation, true, category, context)
		compareField("recommendedInterval", before.RecommendedInterval, after.RecommendedInterval, true, "AoM: recommendedInterval deviation", context)

		compareNumber("SupportedMeasurementInterval", before.SupportedMeasurementInterval, after.SupportedMeasurementInterval, context)
		compareSimpleElements("SupportedMeasurementInterval", before.SupportedMeasurementInterval, after.SupportedMeasurementInterval, context)
	} else if nil != before && nil == after {
		message := fmt.Sprint("SupportedMeasurementIntervals is not exported")
		logDeviation(category, message, context)
	} else if nil == before && nil != after {
		message := fmt.Sprint("SupportedMeasurementIntervals is added in exported PMB")
		logDeviation(category, message, context)
	}
}

func compareField(fieldInfo string, beforeValue string, afterValue string, isOut bool, category string, context Context) {
	if strings.TrimSpace(beforeValue) != strings.TrimSpace(afterValue) {
		if isOut {
			if fieldInfo == "Measurement [UB6FR] networkElementMeasurementId" {
				fmt.Println("======================> ", beforeValue, afterValue)
			}
			message := fmt.Sprint(fieldInfo, " is different, before {", beforeValue, "}, after {", afterValue, "}")
			logDeviation(category, message, context)
		}
	}
}

func raiseDeviation(category string, message string, context Context) {
	deviation := Deviation {
		Level: "error",
		DeviationCategory: category,
		Message: fmt.Sprint(message),
	}
	if context.PMClass != "" {
		deviation.PmClass = &context.PMClass
	}
	if context.MeasurementType != "" {
		deviation.MeasurementType = &context.MeasurementType
	}
	if context.Counter != "" {
		deviation.Counter = &context.Counter
	}
	if context.Dimension != "" {
		deviation.Dimension = &context.Dimension
	}
	deviations = append(deviations, deviation)
}

func log(message string) {
	fmt.Println(message)
	logFile.WriteString(fmt.Sprint(message, "\n"))
}

func logDeviation(category string, message string, context Context) {
	log(message)
	raiseDeviation(category, message, context)
}

func writeResult() {
	output, err := xml.MarshalIndent(result, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	resultFile.Write(output)
}


