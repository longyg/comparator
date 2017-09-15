package main

import (
	"flag"
	"comparator/pm/util"
)

func main()  {
	beforePmb := flag.String("b","", "PMB file path before NIDD")
	afterPmb := flag.String("a", "", "PMB file path after NIDD")
	logPath := flag.String("l", "", "Log file path")
	resultPath := flag.String("r", "", "Result file path")

	flag.Parse()

	pmutil.Compare(*beforePmb, *afterPmb, *logPath, *resultPath)
}