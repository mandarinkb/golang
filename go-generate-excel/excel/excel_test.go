package excel

import "testing"

func TestXxx(t *testing.T) {
	var listExcelData [][]interface{}
	listExcelData = append(listExcelData, []interface{}{
		"1",
		"16/09/23",
		"09:00:00",
		"99",
		"example detail",
	})
	GenerateReportExcelFile("20230916", listExcelData, "../report_template/ex_report.xlsx", "./gen_%s.xlsx")
}
