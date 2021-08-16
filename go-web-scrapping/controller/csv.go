package controller

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
	"webscrapping/service"
)

func WriteCSV(products []service.ProductReponse) error {
	// สร้าง csv file โดยชื่อนั้นตั้งจากวันที่
	fileName := time.Now().Format(time.RFC3339) + ".csv"
	fileName = strings.ReplaceAll(fileName, ":", ".")
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	csvwriter := csv.NewWriter(csvFile)
	var header = []string{("Timestamp"), ("WebName"), ("ProductName"), ("Category"),
		("Price"), ("OriginalPrice"), ("ProductUrl"), ("Image"), ("Icon")}
	for i := range products {
		var row []string
		if i == 0 {
			csvwriter.Write(header) // ใส่ header csv
			row = append(row, products[i].Timestamp)
			row = append(row, products[i].WebName)
			row = append(row, products[i].ProductName)
			row = append(row, products[i].Category)
			row = append(row, products[i].Price)
			row = append(row, products[i].OriginalPrice)
			row = append(row, products[i].ProductUrl)
			row = append(row, products[i].Image)
			row = append(row, products[i].Icon)
			csvwriter.Write(row)
		} else {
			row = append(row, products[i].Timestamp)
			row = append(row, products[i].WebName)
			row = append(row, products[i].ProductName)
			row = append(row, products[i].Category)
			row = append(row, products[i].Price)
			row = append(row, products[i].OriginalPrice)
			row = append(row, products[i].ProductUrl)
			row = append(row, products[i].Image)
			row = append(row, products[i].Icon)
			csvwriter.Write(row)
		}

	}
	csvwriter.Flush()
	csvFile.Close()
	fmt.Println("===== write csv success =====")

	return nil
}

func WriteCSV3Web(products []service.ProductReponse, products2 []service.ProductReponse, products3 []service.ProductReponse) error {
	// สร้าง csv file โดยชื่อนั้นตั้งจากวันที่
	// csvFile, err := os.Create("./" + time.Now().Format(time.RFC3339) + ".csv")
	fileName := time.Now().Format(time.RFC3339) + ".csv"
	fileName = strings.ReplaceAll(fileName, ":", ".")
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	csvwriter := csv.NewWriter(csvFile)
	var header = []string{("Timestamp"), ("WebName"), ("ProductName"), ("Category"),
		("Price"), ("OriginalPrice"), ("ProductUrl"), ("Image"), ("Icon")}
	for i := range products {
		var row []string
		if i == 0 {
			csvwriter.Write(header) // ใส่ header csv
			row = append(row, products[i].Timestamp)
			row = append(row, products[i].WebName)
			row = append(row, products[i].ProductName)
			row = append(row, products[i].Category)
			row = append(row, products[i].Price)
			row = append(row, products[i].OriginalPrice)
			row = append(row, products[i].ProductUrl)
			row = append(row, products[i].Image)
			row = append(row, products[i].Icon)
			csvwriter.Write(row)
		} else {
			row = append(row, products[i].Timestamp)
			row = append(row, products[i].WebName)
			row = append(row, products[i].ProductName)
			row = append(row, products[i].Category)
			row = append(row, products[i].Price)
			row = append(row, products[i].OriginalPrice)
			row = append(row, products[i].ProductUrl)
			row = append(row, products[i].Image)
			row = append(row, products[i].Icon)
			csvwriter.Write(row)
		}
	}
	for i := range products2 {
		var row []string
		row = append(row, products2[i].Timestamp)
		row = append(row, products2[i].WebName)
		row = append(row, products2[i].ProductName)
		row = append(row, products2[i].Category)
		row = append(row, products2[i].Price)
		row = append(row, products2[i].OriginalPrice)
		row = append(row, products2[i].ProductUrl)
		row = append(row, products2[i].Image)
		row = append(row, products2[i].Icon)
		csvwriter.Write(row)
	}
	for i := range products3 {
		var row []string
		row = append(row, products3[i].Timestamp)
		row = append(row, products3[i].WebName)
		row = append(row, products3[i].ProductName)
		row = append(row, products3[i].Category)
		row = append(row, products3[i].Price)
		row = append(row, products3[i].OriginalPrice)
		row = append(row, products3[i].ProductUrl)
		row = append(row, products3[i].Image)
		row = append(row, products3[i].Icon)
		csvwriter.Write(row)
	}
	csvwriter.Flush()
	csvFile.Close()
	fmt.Println("===== write csv success =====")

	return nil
}
