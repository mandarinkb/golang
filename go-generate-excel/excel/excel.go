package excel

import (
	"fmt"
	"time"

	"github.com/tealeg/xlsx"
)

func GenerateReportExcelFile(orderDate string, listExcelData [][]interface{}, templatePath, outboundPath string) (string, error) {
	file, err := xlsx.OpenFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("open file %s error: %+v", templatePath, err)
	}
	date, _ := time.Parse("20060102", orderDate)

	setHeaderToFile(date, len(listExcelData), file)

	currentRow := 5
	setDataToFile(listExcelData, file, currentRow)

	path := getReportPath(date, outboundPath)
	err = file.Save(path)
	if err != nil {
		return "", fmt.Errorf("save file %s error: %+v", path, err)
	}
	return path, nil
}

func GenerateReportExcelFileCustomRow(orderDate string, listExcelData [][]interface{}, templatePath, outboundPath string, currentRow int) (string, error) {
	file, err := xlsx.OpenFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("open file %s error: %+v", templatePath, err)
	}
	date, _ := time.Parse("20060102", orderDate)
	setDataToFile(listExcelData, file, currentRow)

	path := getReportPathWithDateFormat(date, outboundPath, "20060102")
	err = file.Save(path)
	if err != nil {
		return "", fmt.Errorf("save file %s error: %+v", path, err)
	}

	return path, nil
}

func setHeaderToFile(date time.Time, nTxn int, file *xlsx.File) {
	sheet1 := file.Sheets[0]
	rowHeader1 := sheet1.Rows[1]
	cellHeader1 := rowHeader1.Cells[1]
	cellHeader1.SetValue(date.Format("20060102"))

	rowHeader2 := sheet1.Rows[2]
	cellHeader2 := rowHeader2.Cells[1]
	cellHeader2.SetValue(fmt.Sprintf("%d transactions", nTxn))
}

func setDataToFile(listExcelData [][]interface{}, file *xlsx.File, currentRow int) {
	for _, list := range listExcelData {
		file.Sheets[0].AddRow()
		newRow := file.Sheets[0].Rows[currentRow]
		addDateToRow(newRow, list)
		currentRow++
	}
}

func addDateToRow(newRow *xlsx.Row, list []interface{}) {
	for i, v := range list {
		var cell *xlsx.Cell
		if i >= len(newRow.Cells) {
			cell = newRow.AddCell()
		} else {
			cell = newRow.Cells[i]
		}

		cell.SetStyle(getXlsStyle())
		cell.SetValue(v)
	}
}

func getReportPath(date time.Time, outbound string) string {
	return fmt.Sprintf(outbound, date.Format("02012006"))
}

func getReportPathWithDateFormat(date time.Time, outbound, dateFormat string) string {
	return fmt.Sprintf(outbound, date.Format(dateFormat))
}

func addDataToCell(row *xlsx.Row, data string) {
	cell := row.AddCell()
	cell.SetStyle(getXlsStyle()) //check style
	cell.SetValue(data)
}

func getXlsStyle() *xlsx.Style {
	return &xlsx.Style{
		Border: xlsx.Border{
			Left: "thin", LeftColor: "FF000000", Right: "thin", RightColor: "FF000000", Top: "thin", TopColor: "FF000000", Bottom: "thin", BottomColor: "FF000000",
		},
		Fill: xlsx.Fill{
			PatternType: "none",
		},
		Font: xlsx.Font{
			Size: 11, Name: "Calibri", Family: 2, Charset: 0, Color: "FF000000", Bold: false, Italic: false, Underline: false,
		},
		ApplyBorder:    true,
		ApplyFill:      false,
		ApplyFont:      true,
		ApplyAlignment: false,
		Alignment: xlsx.Alignment{
			Indent: 0, ShrinkToFit: false, TextRotation: 0, WrapText: false,
		},
	}
}
