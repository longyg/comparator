package main

import (
	"GoTest/amanual"
)

func main()  {
	beforeDir := "D:\\projects\\Ranger\\IMS\\IMS18.0\\fastpass\\NIDD_test_data\\amanual\\"
	afterDir := "D:\\projects\\Ranger\\IMS\\IMS18.0\\fastpass\\Convertor\\man-converter-17.06\\zip\\adaptation_com.nsn.cscf.man-17.5VI-20170731T100549\\amanual"
	outfile := "D:\\projects\\Ranger\\IMS\\IMS18.0\\fastpass\\NIDD_test_data\\out.txt"

	amanual.Compare(beforeDir, afterDir, outfile)
}
