package main

import (
	"io"
	"log"
	"os"

	"github.com/alexmullins/zip"
)

func main() {
	fzip, err := os.Create("./test.zip")
	if err != nil {
		log.Fatalln(err)
	}
	zipw := zip.NewWriter(fzip)
	defer zipw.Close()

	password := "P@ssw0rd"
	w, err := zipw.Encrypt("your_file_name.csv", password)
	if err != nil {
		log.Fatal(err)
	}

	//read source file
	f1, err := os.Open("./temp/test.csv")
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	//copy source file to zip
	_, err = io.Copy(w, f1)
	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
}
