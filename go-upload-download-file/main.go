package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const FILE_PATH = "data-upload"

// upload file
func uploadFile(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["file"]

	for _, file := range files {
		//path := "data-upload"
		if _, err := os.Stat(FILE_PATH); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(FILE_PATH, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}

		targetPath := filepath.Join("./", FILE_PATH, "/", file.Filename)
		//c.SaveUploadedFile(file, "./"+path+"/"+file.Filename)
		c.SaveUploadedFile(file, targetPath)
	}
	uploadRes := strconv.Itoa(len(files)) + " files uploaded"
	c.IndentedJSON(http.StatusOK, gin.H{"message": uploadRes})
}

func downloadFile(c *gin.Context) {
	fileName := c.Param("file")
	targetPath := filepath.Join(FILE_PATH, "/", fileName)

	//ตรวจสอบ filename เพื่อป้องกันการโจมตีจาก user side
	if !strings.HasPrefix(filepath.Clean(targetPath), FILE_PATH+"/") {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Look like you attacking me"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	//c.File("./data-upload/" + fileName)
	c.File(targetPath)
}
func main() {
	fmt.Println("Listening and serving HTTP on :8080")

	// set release mode
	// using env:   export GIN_MODE=release
	// using code:  gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// CORS gin's middleware Default() allows all origins
	router.Use(cors.Default())

	router.POST("/upload", uploadFile)
	router.Static("/image", "./data-upload") //  จะ return folder ที่เก็บไฟล์ให้ ตอนเรียกอ้างอิงไฟล์ได้เลย http://127.0.0.1:8080/image/demo.png
	router.GET("/download/:file", downloadFile)
	router.Run(":8080")

}
