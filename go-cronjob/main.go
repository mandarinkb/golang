package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CronJob struct {
	CronExpress string `json:"CronExpress"`
	ScriptPath  string `json:"ScriptPath"`
	ScriptName  string `json:"ScriptName"` //unique
	LogPath     string `json:"LogPath"`
}

type Response struct {
	Data []CronJobDataDetail `json:"Data"`
}
type CronJobDataDetail struct {
	CronJobData string `json:"CronJobData"`
}

type ResponseLog struct {
	FileName string
	Message  string
}

// 0 * * * * root . $HOME/.profile; cd /data/jobs/minorcrm-jobs-go/shellscript/jobs_script && sh job_dashboard_card_1016 >> /data/jobs/minorcrm-jobs-go/shellscript/job_dashboard_card_1016.log 2>&1
// var cronTemplate = "%s root . $HOME/.profile; cd %s && sh %s >> %s.log 2>&1\n"
var cronTemplate = "%s %s/%s.sh >> %s/%s.log 2>&1\n"
var path = "../etc/cron.d/root"

func main() {
	// fiber instance
	app := fiber.New()

	// routes
	app.Get("/cron", readJob)
	app.Post("/cron", createJob)
	app.Put("/cron", updateJob)
	app.Delete("/cron", deleteJob)

	// app listening at PORT: 3000
	app.Listen(":3000")
}
func readJob(c *fiber.Ctx) error {
	readFile, err := os.Open(path)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}
	defer readFile.Close()

	response := Response{}
	listCronJobDataDetail := []CronJobDataDetail{}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		cronData := CronJobDataDetail{}
		cronData.CronJobData = fileScanner.Text()
		listCronJobDataDetail = append(listCronJobDataDetail, cronData)
	}

	response.Data = listCronJobDataDetail
	return c.Status(200).JSON(response)
}

func createJob(c *fiber.Ctx) error {
	resp := CronJob{}
	if err := c.BodyParser(&resp); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}

	if isDupcateName(resp.ScriptName) {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": "script name is duplicate",
		})
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}
	defer f.Close()

	data := fmt.Sprintf(cronTemplate, resp.CronExpress, resp.ScriptPath, resp.ScriptName, resp.LogPath, resp.ScriptName)
	_, err = f.WriteString(data)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err,
		})
	}
	reload()
	return c.Status(200).JSON(&fiber.Map{"message": "create done."})
}

func updateJob(c *fiber.Ctx) error {
	resp := CronJob{}
	if err := c.BodyParser(&resp); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}

	readFile, err := os.Open(path)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}
	defer readFile.Close()

	listLine := []string{}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Contains(line, resp.ScriptName) {
			splitLine := strings.Split(line, "/app/script/")
			line = fmt.Sprintf("%s /app/script/%s", resp.CronExpress, splitLine[1])
		}
		listLine = append(listLine, line)
	}
	text := strings.Join(listLine, "\n")
	text = fmt.Sprintf("%s\n", text)
	// / write file
	writeFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}
	defer writeFile.Close()

	_, err = writeFile.WriteString(text)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err,
		})
	}
	reload()
	return c.Status(200).JSON(&fiber.Map{"message": "update done."})
}

func deleteJob(c *fiber.Ctx) error {
	id := c.Query("script_name")

	readFile, err := os.Open(path)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}
	defer readFile.Close()

	listLine := []string{}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Contains(line, id) {
			listLine = append(listLine, line)
		}
	}
	text := strings.Join(listLine, "\n")
	text = fmt.Sprintf("%s\n", text)
	// / write file
	writeFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}
	defer writeFile.Close()

	_, err = writeFile.WriteString(text)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  400,
			"message": err,
		})
	}
	reload()
	return c.Status(200).JSON(&fiber.Map{"message": "delete done."})
}

func isDupcateName(fileName string) bool {
	readFile, err := os.Open(path)
	if err != nil {
		return true
	}
	defer readFile.Close()

	isDup := false
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), fileName) {
			isDup = true
		}
	}

	return isDup
}
func reload() {
	// cmd := exec.Command("bash", "-c", " cd /app/script/ && sh job_reload.sh")
	cmd := exec.Command("bash", "-c", " crontab /etc/cron.d/root")
	err := cmd.Run()
	if err != nil {
		fmt.Println("error: ", err)
	}
}

// validate script name not duplicate name
// read data in file
// edit data in file
// delete data in file
// write cmd in file sh
// create docker file
// สามารถ วาง root file in to cron job และ CRUD ได้
// read log file เมื่อ job ทำงาน
