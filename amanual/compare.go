package amanual

import (
	"fmt"
	"os"
	"strings"
)

var file *os.File

func Compare(beforeDir string, afterDir string, outfile string) {
	beforeMans, afterMans := readMans(beforeDir, afterDir)

	file, _ = os.Create(outfile)
	defer file.Close()

	compareNumbers(beforeMans, afterMans)

	compareContent(beforeMans, afterMans)
}

func readMans(beforeDir string, afterDir string) (map[string]Manual, map[string]Manual) {
	beforeMans := ReadMans(beforeDir)
	afterMans := ReadMans(afterDir)
	return beforeMans, afterMans
}

func compareNumbers(beforeMans map[string]Manual, afterMans map[string]Manual) {
	if len(beforeMans) == len(afterMans) {
		fmt.Println("Number of man pages are same:", len(beforeMans))
		write(fmt.Sprint("Number of man pages are same: ", len(beforeMans)))
	} else {
		fmt.Println("Number of man pages are different, before:", len(beforeMans), ", after:", len(afterMans))
		write(fmt.Sprint("Number of man pages are different, before: ", len(beforeMans), ", after: ", len(afterMans)))

		for beforeName, _ := range beforeMans {
			isFound := false
			for afterName, _ := range afterMans {
				if beforeName == afterName {
					isFound = true
					break
				}
			}
			if !isFound {
				fmt.Println(beforeName, "is missing after NIDD")
				write(fmt.Sprint(beforeName, " is missing after NIDD"))
			}
		}

		for afterName, _ := range afterMans {
			isFound := false
			for beforeName, _ := range beforeMans {
				if afterName == beforeName {
					isFound = true
					break
				}
			}
			if !isFound {
				fmt.Println(afterName, "is adding after NIDD")
				write(fmt.Sprint(afterName, " is adding after NIDD"))
			}
		}
	}
}

func compareContent(beforeMans map[string]Manual, afterMans map[string]Manual) {
	for beforeName, beforeMan := range beforeMans {
		for afterName, afterMan := range afterMans {
			if beforeName == afterName {
				compareManContent(beforeName, beforeMan, afterMan)
			}
		}

	}
}

func compareManContent(manName string, beforeMan Manual, afterMan Manual) {
	compareField("alarmText", beforeMan.AlarmText, afterMan.AlarmText, manName, false)
	compareField("specificProblem", beforeMan.SpecificProblem, afterMan.SpecificProblem, manName, false)
	compareField("adaptationHref", beforeMan.Adaptation.Href, afterMan.Adaptation.Href, manName, false)
	compareField("xmi", beforeMan.XMIXmlns, afterMan.XMIXmlns, manName, false)
	compareField("patchLevel", beforeMan.PatchLevel, afterMan.PatchLevel, manName, false)
	compareField("alarmType", beforeMan.AlarmType, afterMan.AlarmType, manName, false)
	compareField("cancelling", beforeMan.Cancelling, afterMan.Cancelling, manName, false)
	compareField("instructions", beforeMan.Instructions, afterMan.Instructions, manName, false)
	compareField("meaning", beforeMan.Meaning, afterMan.Meaning, manName, false)
	compareField("perceivedSeverityInfo", beforeMan.PerceivedSeverityInfo, afterMan.PerceivedSeverityInfo, manName, false)
	compareField("probableCause", beforeMan.ProbableCause, afterMan.ProbableCause, manName, false)
	compareField("supplementaryInformationFields", beforeMan.SupplementaryInformationFields, afterMan.SupplementaryInformationFields, manName, false)
	compareField("XMIVersion", beforeMan.XMIVersion, afterMan.XMIVersion, manName, true)
	compareField("xmlns", beforeMan.Xmlns, afterMan.Xmlns, manName, false)

}

func compareField(fieldName string, beforeValue string, afterValue string, manName string, isOut bool) {
	if strings.TrimSpace(beforeValue) != strings.TrimSpace(afterValue) {
		if isOut {
			fmt.Println("==================== Compare man page:", manName, "=======================")
			fmt.Println(fieldName, "is different:")
			fmt.Println("	before : {",beforeValue,"}")
			fmt.Println("	after  : {",afterValue,"}")

			write(fmt.Sprint("==================== Compare man page:", manName, "======================="))
			write(fmt.Sprint(fieldName, " is different:"))
			write(fmt.Sprint("	before : {",beforeValue,"}"))
			write(fmt.Sprint("	after  : {",afterValue,"}"))
		}
	}
}

func write(message string) {
	file.WriteString(fmt.Sprint(message,"\n"))
}