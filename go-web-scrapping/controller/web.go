package controller

import (
	"fmt"
	"time"
	"webscrapping/repository"
	"webscrapping/service"
)

func Bigc() {
	start := time.Now().Format(time.RFC3339)
	bigc, err := repository.NewBigC().Bigc()
	if err != nil {
		fmt.Println(err)
	}
	bigcData, err := service.NewProductService().CsvData(bigc)
	if err != nil {
		fmt.Println(err)
	}
	WriteCSV(bigcData)
	end := time.Now().Format(time.RFC3339)
	fmt.Println("start : ", start)
	fmt.Println("end   : ", end)
}
func Makro() {
	start := time.Now().Format(time.RFC3339)
	makro, err := repository.NewMakro().Makro()
	if err != nil {
		fmt.Println(err)
	}
	makroData, err := service.NewProductService().CsvData(makro)
	if err != nil {
		fmt.Println(err)
	}
	WriteCSV(makroData)
	end := time.Now().Format(time.RFC3339)
	fmt.Println("start : ", start)
	fmt.Println("end   : ", end)
}
func Lotus() {
	start := time.Now().Format(time.RFC3339)
	lotus, err := repository.NewLotus().Lotus()
	if err != nil {
		fmt.Println(err)
	}
	lotusData, err := service.NewProductService().CsvData(lotus)
	if err != nil {
		fmt.Println(err)
	}
	WriteCSV(lotusData)
	end := time.Now().Format(time.RFC3339)
	fmt.Println("start : ", start)
	fmt.Println("end   : ", end)
}
func All() {
	start := time.Now().Format(time.RFC3339)
	// Bigc
	bigc, err := repository.NewBigC().Bigc()
	if err != nil {
		fmt.Println(err)
	}
	bigcData, err := service.NewProductService().CsvData(bigc)
	if err != nil {
		fmt.Println(err)
	}
	// Makro
	makro, err := repository.NewMakro().Makro()
	if err != nil {
		fmt.Println(err)
	}
	makroData, err := service.NewProductService().CsvData(makro)
	if err != nil {
		fmt.Println(err)
	}
	// Lotus
	lotus, err := repository.NewLotus().Lotus()
	if err != nil {
		fmt.Println(err)
	}
	lotusData, err := service.NewProductService().CsvData(lotus)
	if err != nil {
		fmt.Println(err)
	}
	// controller.WriteCSV(pds)
	WriteCSV3Web(bigcData, makroData, lotusData)
	end := time.Now().Format(time.RFC3339)
	fmt.Println("start : ", start)
	fmt.Println("end   : ", end)
}
