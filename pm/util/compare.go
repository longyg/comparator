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

func compareBasicInfo(beforePmb PMBasic, afterPmb PMBasic) {
	fmt.Println("==================== Compare PMB basic information =======================")
	write(fmt.Sprint("==================== Compare PMB basic information ======================="))
	devType := "PMB Basic Info Deviation"
	compareField("xmlns:ns2", beforePmb.Xmlns2, afterPmb.Xmlns2, true, devType, "", "", "")
	compareField("xmlns:ns3", beforePmb.Xmlns3, afterPmb.Xmlns3, true, devType, "", "", "")
	compareField("xmlns:xsi", beforePmb.XmlnsXSI, afterPmb.XmlnsXSI, true, devType, "", "", "")
	compareField("schemaVersion", beforePmb.SchemaVersion, afterPmb.SchemaVersion, true, devType, "", "", "")
	compareField("interfaceVersion", beforePmb.InterfaceVersion, afterPmb.InterfaceVersion, true, devType, "", "", "")
	compareField("schemaLocation", beforePmb.SchemaLocation, afterPmb.SchemaLocation, true, devType, "", "", "")
	compareField("Adaptation.href", beforePmb.Adaptation.Href, afterPmb.Adaptation.Href, true, devType, "", "", "")
	compareField("Product.typeName", beforePmb.Product.TypeName, afterPmb.Product.TypeName, true, devType, "", "", "")
	compareField("Product.vendor", beforePmb.Product.Vendor, afterPmb.Product.Vendor, true, devType, "", "", "")
	compareField("Product.release", beforePmb.Product.Release, afterPmb.Product.Release, true, devType, "", "", "")
	compareField("Product.Description", beforePmb.Product.Description, afterPmb.Product.Description, true, devType, "", "", "")
}

func compareElementNumber(element string, before interface{}, after interface{})  {
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
		fmt.Println("Number of", element, "is not equal, before:", lenBefore, ", after:", lenAfter)
		write(fmt.Sprint("Number of", element, "is not equal, before:", lenBefore, ", after:", lenAfter))
		deviation := Deviation {
			Level: "error",
			Type: fmt.Sprint(element, "Number Deviation"),
			Message: fmt.Sprint("Number of", element, "is not equal, before:", lenBefore, ", after:", lenAfter),
		}
		deviations = append(deviations, deviation)
	}
}

func compareElements(before interface{}, after interface{}, fieldName string) {
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
		bTypeName := be.Type().Name()
		bFieldValue := be.FieldByName(fieldName).String()

		isFound := false

		for j := 0; j < av.Len(); j++ {
			ae := av.Index(j)
			aFieldValue := ae.FieldByName(fieldName).String()

			if bFieldValue == aFieldValue {
				isFound = true
				compareElement(be, ae)
				break
			}
		}
		if !isFound {
			fmt.Println(bTypeName, "[", bFieldValue, "] is not found in exported PMB")
			write(fmt.Sprint(bTypeName, "[", bFieldValue, "] is not found in exported PMB"))
			deviation := Deviation {
				Level: "error",
				Type: "Missing PMClassInfo",
				PmClass: bFieldValue,
				Message: fmt.Sprint(fmt.Sprint(bTypeName, "[", bFieldValue, "] is not found in exported PMB")),
			}
			deviations = append(deviations, deviation)
		}
	}
}

func compareElement(before interface{}, after interface{}) {

}

func comparePMClasses(beforePmb PMBasic, afterPmb PMBasic)  {
	compareElementNumber("PMClassInfo", beforePmb.PMClasses.PMClassInfo, afterPmb.PMClasses.PMClassInfo)

	compareElements(beforePmb.PMClasses.PMClassInfo, afterPmb.PMClasses.PMClassInfo, "Name")

	for _, beforeClass := range beforePmb.PMClasses.PMClassInfo {
		isFound := false
		for _, afterClass := range afterPmb.PMClasses.PMClassInfo {
			if beforeClass.Name == afterClass.Name {
				isFound = true
				comparePMClassInfo(beforeClass, afterClass)
				break
			}
		}
		if !isFound {
			fmt.Println("PM class:", beforeClass.Name, "is not found in exported PMB")
			write(fmt.Sprint("PM class:", beforeClass.Name, "is not found in exported PMB"))
			deviation := Deviation {
				Level: "error",
				Type: "Missing PMClassInfo",
				PmClass: beforeClass.Name,
				Message: fmt.Sprint(fmt.Sprint("PM class:", beforeClass.Name, "is not found in exported PMB")),
			}
			deviations = append(deviations, deviation)
		}
	}

	for _, afterClass := range afterPmb.PMClasses.PMClassInfo {
		isFound := false
		for _, beforeClass := range beforePmb.PMClasses.PMClassInfo {
			if afterClass.Name == beforeClass.Name {
				isFound = true
				break
			}
		}
		if !isFound {
			fmt.Println("PM class:", afterClass.Name, "is added wrongly in exported PMB")
			write(fmt.Sprint("PM class:", afterClass.Name, "is added wrongly in exported PMB"))
			deviation := Deviation {
				Level: "error",
				Type: "Adding PMClassInfo",
				PmClass: afterClass.Name,
				Message: fmt.Sprint(fmt.Sprint("PM class:", afterClass.Name, "is added wrongly in exported PMB")),
			}
			deviations = append(deviations, deviation)
		}
	}
}

func comparePMClassInfo(beforeClass PMClassInfo, afterClass PMClassInfo) {
	fmt.Println("==================== Compare PMClassInfo:", beforeClass.Name, "=======================")
	write(fmt.Sprint("==================== Compare PMClassInfo:", beforeClass.Name, "======================="))
	devType := "PMClassInfo Deviation"
	compareField("PMClassInfo [" + beforeClass.Name + "] nameInOmes", beforeClass.NameInOMeS, afterClass.NameInOMeS, true, devType, beforeClass.Name, "", "")
	compareField("[" + beforeClass.Name + "] intransient", beforeClass.Intransient, afterClass.Intransient, true, devType, beforeClass.Name, "", "")
	compareField("[" + beforeClass.Name + "] presentation", beforeClass.Presentation, afterClass.Presentation, true, devType, beforeClass.Name, "", "")
	compareField("[" + beforeClass.Name + "] Description", beforeClass.Description, afterClass.Description, true, devType, beforeClass.Name, "", "")
}

func compareAdminOfMeasurements(before *AdminOfMeasurements, after *AdminOfMeasurements)  {
	fmt.Println("==================== Compare AdminOfMeasurements =======================")
	write(fmt.Sprint("==================== Compare AdminOfMeasurements ======================="))
	if nil != before && nil != after {
		compareSupportedMeasurementIntervals(before.SupportedMeasurementIntervals, after.SupportedMeasurementIntervals, Measurement{})
	} else if nil != before && nil == after {
		fmt.Println("AdminOfMeasurements is not exported")
		write(fmt.Sprint("AdminOfMeasurements is not exported"))
	} else if nil == before && nil != after {
		fmt.Println("AdminOfMeasurements is added in exported PMB")
		write(fmt.Sprint("AdminOfMeasurements is added in exported PMB"))
	} else {
		fmt.Println("AdminOfMeasurements is not defined either before NIDD and after NIDD")
		write(fmt.Sprint("AdminOfMeasurements is not defined either before NIDD and after NIDD"))
	}
}

func compareMeasurements(beforePmb PMBasic, afterPmb PMBasic) {
	lenBefore := len(beforePmb.Measurement)
	lenAfter := len(afterPmb.Measurement)
	if lenBefore != lenAfter {
		fmt.Println("Number of Measurement is not equal, before:", lenBefore, ", after:", lenAfter)
		write(fmt.Sprint("Number of Measurement is not equal, before:", lenBefore, ", after:", lenAfter))
	} else {
		for _, before := range beforePmb.Measurement {
			isFound := false
			for _, after := range afterPmb.Measurement {
				if before.MeasurementType == after.MeasurementType {
					isFound = true
					compareMeasurement(before, after)
					break
				}
			}
			if !isFound {
				fmt.Println("PM class:", before.MeasurementType, "is not found in exported PMB")
				write(fmt.Sprint("PM class:", before.MeasurementType, "is not found in exported PMB"))
			}
		}

		for _, after := range afterPmb.Measurement {
			isFound := false
			for _, before := range beforePmb.Measurement {
				if after.MeasurementType == before.MeasurementType {
					isFound = true
					break
				}
			}
			if !isFound {
				fmt.Println("PM class:", after.MeasurementType, "is added wrongly in exported PMB")
				write(fmt.Sprint("PM class:", after.MeasurementType, "is added wrongly in exported PMB"))
			}
		}
	}
}

func compareMeasurement(before Measurement, after Measurement) {
	fmt.Println("==================== Compare measurement:", before.MeasurementType, "=======================")
	write(fmt.Sprint("==================== Compare measurement:", before.MeasurementType, "======================="))
	compareMeasurementField(before, after)
	compareMeasuredTargets(before, after)
	compareCounters(before, after)
	compareSupportedMeasurementIntervals(before.SupportedMeasurementIntervals, after.SupportedMeasurementIntervals, before)
}

func compareCounters(beforeMeas Measurement, afterMeas Measurement)  {
	lenBefore := len(beforeMeas.MeasuredIndicator)
	lenAfter := len(afterMeas.MeasuredIndicator)
	if lenBefore != lenAfter {
		fmt.Println("Number of Counters is not equal, before:", lenBefore, ", after:", lenAfter)
		write(fmt.Sprint("Number of Counters is not equal, before:", lenBefore, ", after:", lenAfter))
	} else {
		for _, before := range beforeMeas.MeasuredIndicator {
			isFound := false
			for _, after := range afterMeas.MeasuredIndicator {
				if before.Name == after.Name {
					isFound = true
					compareCounter(before, after, beforeMeas)
					break
				}
			}
			if !isFound {
				fmt.Println("Counter:", before.Name, "is not found in exported PMB")
				write(fmt.Sprint("Counter:", before.Name, "is not found in exported PMB"))
			}
		}

		for _, after := range afterMeas.MeasuredIndicator {
			isFound := false
			for _, before := range beforeMeas.MeasuredIndicator {
				if after.Name == before.Name {
					isFound = true
					break
				}
			}
			if !isFound {
				fmt.Println("Counter:", after.Name, "is added wrongly in exported PMB")
				write(fmt.Sprint("Counter:", after.Name, "is added wrongly in exported PMB"))
			}
		}
	}
}

func compareCounter(before MeasuredIndicator, after MeasuredIndicator, meas Measurement)  {
	compareCounterFields(before, after, meas)
	compareAggRules(before, after, meas)
	compareFormula(before.TimeAndObjectAggregationFormula, after.TimeAndObjectAggregationFormula)
	compareDocumentation(before.Documentation, after.Documentation)
	compareSupportedInProducts(before.SupportedInProducts, after.SupportedInProducts)
}

func compareAggRules(before MeasuredIndicator, after MeasuredIndicator, meas Measurement)  {
	devType := "Counter AGG Rules Deviation"
	if "" != before.TimeAndObjectAggregationRule && "" != after.TimeAndObjectAggregationRule {
		compareField("Counter [" + before.Name + "] TimeAndObjectAggregationRule", before.TimeAndObjectAggregationRule, after.TimeAndObjectAggregationRule, true, devType, "", meas.MeasurementType, before.Name)
	} else if "" != before.TimeAndObjectAggregationRule && "" == after.TimeAndObjectAggregationRule {
		if "" != after.TimeAggregationRule && "" != after.ObjectAggregationRule {
			compareField("Counter [" + before.Name + "] TimeAggregationRule", before.TimeAndObjectAggregationRule, after.TimeAggregationRule, true, devType, "", meas.MeasurementType, before.Name)
			compareField("Counter [" + before.Name + "] ObjectAggregationRule", before.TimeAndObjectAggregationRule, after.ObjectAggregationRule, true, devType, "", meas.MeasurementType, before.Name)
		} else {
			devType = "Counter AGG Rule Not Exported"
			if "" == after.TimeAggregationRule {
				fmt.Println("TimeAggregationRule is not defined in exported PMB")
				write(fmt.Sprint("TimeAggregationRule is not defined in exported PMB"))
				deviation := Deviation{
					Level: "error",
					MeasurementType: meas.MeasurementType,
					Counter: before.Name,
					Type: devType,
					Message: fmt.Sprint("TimeAggregationRule is not defined in exported PMB"),
				}
				deviations = append(deviations, deviation)
			}
			if "" == after.ObjectAggregationRule {
				fmt.Println("ObjectAggregationRule is not defined")
				write(fmt.Sprint("ObjectAggregationRule is not defined in exported PMB"))
				deviation := Deviation{
					Level: "error",
					MeasurementType: meas.MeasurementType,
					Counter: before.Name,
					Type: devType,
					Message: fmt.Sprint("ObjectAggregationRule is not defined in exported PMB"),
				}
				deviations = append(deviations, deviation)
			}
		}
	} else if "" == before.TimeAndObjectAggregationRule && "" != after.TimeAndObjectAggregationRule {
		if "" != before.TimeAggregationRule && "" != before.ObjectAggregationRule {
			compareField("Counter [" + before.Name + "] TimeAggregationRule", before.TimeAggregationRule, after.TimeAndObjectAggregationRule, true, devType, "", meas.MeasurementType, before.Name)
			compareField("Counter [" + before.Name + "] ObjectAggregationRule", before.ObjectAggregationRule, after.TimeAndObjectAggregationRule, true, devType, "", meas.MeasurementType, before.Name)
		} else {
			devType = "Counter AGG Rule Not Defined"
			if "" == before.TimeAggregationRule {
				fmt.Println("TimeAggregationRule is not defined in import PMB")
				write(fmt.Sprint("TimeAggregationRule is not defined in import PMB"))
				deviation := Deviation{
					Level: "error",
					MeasurementType: meas.MeasurementType,
					Counter: before.Name,
					Type: devType,
					Message: fmt.Sprint("TimeAggregationRule is not defined in import PMB"),
				}
				deviations = append(deviations, deviation)
			}
			if "" == before.ObjectAggregationRule {
				fmt.Println("ObjectAggregationRule is not defined in import PMB")
				write(fmt.Sprint("ObjectAggregationRule is not defined in import PMB"))
				deviation := Deviation{
					Level: "error",
					MeasurementType: meas.MeasurementType,
					Counter: before.Name,
					Type: devType,
					Message: fmt.Sprint("ObjectAggregationRule is not defined in import PMB"),
				}
				deviations = append(deviations, deviation)
			}
		}
	} else {
		compareField("Counter [" + before.Name + "] TimeAggregationRule", before.TimeAggregationRule, after.TimeAggregationRule, true, devType, "", meas.MeasurementType, before.Name)
		compareField("Counter [" + before.Name + "] ObjectAggregationRule", before.ObjectAggregationRule, after.ObjectAggregationRule, true, devType, "", meas.MeasurementType, before.Name)
	}
}

func compareSupportedInProducts(before SupportedInProducts, after SupportedInProducts)  {

}

func compareDocumentation(before MeasuredIndicatorDocumentation, after MeasuredIndicatorDocumentation)  {

}

func compareFormula(before Formula, after Formula)  {

}

func compareCounterFields(before MeasuredIndicator, after MeasuredIndicator, meas Measurement)  {
	b := before.NameInOMeS
	a := after.NameInOMeS
	if "" == b {
		b = before.Name
	}
	if "" == a {
		a = after.Name
	}
	devType := "Counter Deviation"
	compareField("Counter [" + before.Name + "] nameInOMeS", b, a, true, devType, "", meas.MeasurementType, before.Name)

	compareField("Counter [" + before.Name + "] Description", before.Description, after.Description, true, devType, "", meas.MeasurementType, before.Name)
	compareField("Counter [" + before.Name + "] dataCollectionMethod", before.DataCollectionMethod, after.DataCollectionMethod, true, devType, "", meas.MeasurementType, before.Name)
	b = before.Presentation
	a = after.Presentation
	if "" == b {
		b = before.Name
	}
	if "" == a {
		a = after.Name
	}
	compareField("Counter [" + before.Name + "] presentation", b, a, true, devType, "", meas.MeasurementType, before.Name)
	compareField("Counter [" + before.Name + "] unit", before.Unit, after.Unit, true, devType, "", meas.MeasurementType, before.Name)
	compareField("Counter [" + before.Name + "] defaultValue", before.DefaultValue, after.DefaultValue, true, devType, "", meas.MeasurementType, before.Name)
	compareField("Counter [" + before.Name + "] minValue", before.MinValue, after.MinValue, true, devType, "", meas.MeasurementType, before.Name)
	compareField("Counter [" + before.Name + "] maxValue", before.MaxValue, after.MaxValue, true, devType, "", meas.MeasurementType, before.Name)
}

func compareSupportedMeasurementIntervals(before *SupportedMeasurementIntervals, after *SupportedMeasurementIntervals, meas Measurement)  {
	devType := "AoM Config Deviation"
	if nil != before && nil != after {
		compareField("Documentation", before.Documentation, after.Documentation, true, devType, "", meas.MeasurementType, "")
		compareField("recommendedInterval", before.RecommendedInterval, after.RecommendedInterval, true, "AoM: recommendedInterval deviation", "", meas.MeasurementType, "")
		lenBefore := len(before.SupportedMeasurementInterval)
		lenAfter := len(after.SupportedMeasurementInterval)
		if lenBefore != lenAfter {
			fmt.Println("Number of SupportedMeasurementInterval is not equal, before:", lenBefore, ", after:", lenAfter)
			write(fmt.Sprint("Number of SupportedMeasurementInterval is not equal, before:", lenBefore, ", after:", lenAfter))
		} else {
			for _, before := range before.SupportedMeasurementInterval {
				isFound := false
				for _, after := range after.SupportedMeasurementInterval {
					if before == after {
						isFound = true
						break
					}
				}
				if !isFound {
					fmt.Println("SupportedMeasurementInterval:", before, "is not found in exported PMB")
					write(fmt.Sprint("SupportedMeasurementInterval:", before, "is not found in exported PMB"))
				}
			}

			for _, after := range after.SupportedMeasurementInterval {
				isFound := false
				for _, before := range before.SupportedMeasurementInterval {
					if after == before {
						isFound = true
						break
					}
				}
				if !isFound {
					fmt.Println("SupportedMeasurementInterval:", after, "is added wrongly in exported PMB")
					write(fmt.Sprint("SupportedMeasurementInterval:", after, "is added wrongly in exported PMB"))
				}
			}
		}
	} else if nil != before && nil == after {
		fmt.Println("SupportedMeasurementIntervals is not exported")
		write(fmt.Sprint("SupportedMeasurementIntervals is not exported"))
	} else if nil == before && nil != after {
		fmt.Println("SupportedMeasurementIntervals is added in exported PMB")
		write(fmt.Sprint("SupportedMeasurementIntervals is added in exported PMB"))
	}
}

func compareMeasuredTargets(before Measurement, after Measurement) {
	lenBefore := len(before.MeasuredTarget)
	lenAfter := len(after.MeasuredTarget)
	if lenBefore != lenAfter {
		fmt.Println("Number of MeasuredTarget is not equal, before:", lenBefore, ", after:", lenAfter)
		write(fmt.Sprint("Number of MeasuredTarget is not equal, before:", lenBefore, ", after:", lenAfter))
	} else {
		for _, before := range before.MeasuredTarget {
			isFound := false
			for _, after := range after.MeasuredTarget {
				if before.Dimension == after.Dimension {
					isFound = true
					compareHierarchy(before, after)
					break
				}
			}
			if !isFound {
				fmt.Println("MeasuredTarget with diemension:", before.Dimension, "is not found in exported PMB")
				write(fmt.Sprint("MeasuredTarget with diemension:", before.Dimension, "is not found in exported PMB"))
			}
		}

		for _, after := range after.MeasuredTarget {
			isFound := false
			for _, before := range before.MeasuredTarget {
				if after.Dimension == before.Dimension {
					isFound = true
					break
				}
			}
			if !isFound {
				fmt.Println("PM class:", after.Dimension, "is added wrongly in exported PMB")
				write(fmt.Sprint("PM class:", after.Dimension, "is added wrongly in exported PMB"))
			}
		}
	}
}

func compareHierarchy(before MeasuredTarget, after MeasuredTarget) {
	lenBefore := len(before.Hierarchy)
	lenAfter := len(after.Hierarchy)
	if lenBefore != lenAfter {
		fmt.Println("Number of Hierarchy is not equal, before:", lenBefore, ", after:", lenAfter)
		write(fmt.Sprint("Number of Hierarchy is not equal, before:", lenBefore, ", after:", lenAfter))
	} else {
		for _, before := range before.Hierarchy {
			isFound := false
			for _, after := range after.Hierarchy {
				if before == after {
					isFound = true
					break
				}
			}
			if !isFound {
				fmt.Println("Hierarchy:", before, "is not found in exported PMB")
				write(fmt.Sprint("Hierarchy:", before, "is not found in exported PMB"))
			}
		}

		for _, after := range after.Hierarchy {
			isFound := false
			for _, before := range before.Hierarchy {
				if after == before {
					isFound = true
					break
				}
			}
			if !isFound {
				fmt.Println("Hierarchy:", after, "is added wrongly in exported PMB")
				write(fmt.Sprint("Hierarchy:", after, "is added wrongly in exported PMB"))
			}
		}
	}
}

func compareMeasurementField(before Measurement, after Measurement) {
	b := before.MeasurementTypeInOMeS
	a := after.MeasurementTypeInOMeS
	if "" == b {
		b = before.MeasurementType
	}
	if "" == a {
		a = after.MeasurementType
	}
	devType := "Measurement Deviation"
	compareField("Measurement [" + before.MeasurementType + "] measurementTypeInOMeS", b, a, true, devType, "", before.MeasurementType, "")
	b = before.Presentation
	a = after.Presentation
	if "" == b {
		b = before.MeasurementType
	}
	if "" == a {
		a = after.MeasurementType
	}
	compareField("Measurement [" + before.MeasurementType + "] presentation", b, a, true, devType, "", before.MeasurementType, "")
	compareField("Measurement [" + before.MeasurementType + "] defaultInterval", before.DefaultInterval, after.DefaultInterval, true, devType, "", before.MeasurementType, "")
	compareField("Measurement [" + before.MeasurementType + "] networkElementMeasurementId", before.NetworkElementMeasurementId, after.NetworkElementMeasurementId, true, "Measurement: networkElementMeasurementId deviation", "", before.MeasurementType, "")
	compareField("Measurement [" + before.MeasurementType + "] omesFileGroupName", before.OmesFileGroupName, after.OmesFileGroupName, true, devType, "", before.MeasurementType, "")
	compareField("Measurement [" + before.MeasurementType + "] resultsPerInterval", before.ResultsPerInterval, after.ResultsPerInterval, true, devType, "", before.MeasurementType, "")
	b = before.AomSupported
	a = after.AomSupported
	if "" == b {
		b = "true"
	}
	if "" == a {
		a = "true"
	}
	compareField("[" + before.MeasurementType + "] aomSupported", b, a, true, devType, "", before.MeasurementType, "")
	compareField("[" + before.MeasurementType + "] Description", before.Description, after.Description, true, devType, "", before.MeasurementType, "")
}

func write(message string) {
	logFile.WriteString(fmt.Sprint(message,"\n"))
}

func compareField(fieldName string, beforeValue string, afterValue string, isOut bool, deviationType string, pmClass string, measurementType string, counter string) {
	if strings.TrimSpace(beforeValue) != strings.TrimSpace(afterValue) {
		if isOut {
			fmt.Println(fieldName, "is different:")
			fmt.Println("	before : {",beforeValue,"}")
			fmt.Println("	after  : {",afterValue,"}")

			write(fmt.Sprint(fieldName, " is different:"))
			write(fmt.Sprint("	before : {",beforeValue,"}"))
			write(fmt.Sprint("	after  : {",afterValue,"}"))

			deviation := Deviation{
				Level: "error",
				Type: deviationType,
				PmClass: pmClass,
				MeasurementType: measurementType,
				Counter: counter,
				Message: fmt.Sprint(fieldName, " is different: before {", beforeValue, "}, after {", afterValue, "}"),
			}
			deviations = append(deviations, deviation)
		}
	}
}

func writeResult() {
	output, err := xml.MarshalIndent(result, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	resultFile.Write(output)
}


