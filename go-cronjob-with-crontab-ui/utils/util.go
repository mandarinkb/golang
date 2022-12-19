package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-cronjob-with-crontab-ui/logger"
	"github.com/go-cronjob-with-crontab-ui/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func RemoveElement(s []model.CronJobData, i int) []model.CronJobData {
	s[i] = s[len(s)-1]
	s = s[:len(s)-1]
	return s
}

func ReadFile(file string) ([]model.Crontab, error) {
	var fileLines []model.Crontab
	readFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		crontab := model.Crontab{}

		err := json.Unmarshal([]byte(fileScanner.Text()), &crontab)
		if err != nil {
			return nil, err
		}

		fileLines = append(fileLines, crontab)
	}
	return fileLines, nil

}

func ReadFileSystemEventChange(file string, callback func()) {
	//filename := "./crontab/crontab-ui-data/crontab.db"
	err := waitUntilFind(file)
	if err != nil {
		log.Fatalln(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	err = watcher.Add(file)
	if err != nil {
		log.Fatalln(err)
	}

	renameCh := make(chan bool)
	removeCh := make(chan bool)
	errCh := make(chan error)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					// log.Printf("Write:  %s: %s", event.Op, event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					// log.Printf("Create: %s: %s", event.Op, event.Name)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					// log.Printf("Remove: %s: %s", event.Op, event.Name)
					removeCh <- true
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					// log.Printf("Rename: %s: %s", event.Op, event.Name)
					renameCh <- true
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Printf("Chmod:  %s: %s", event.Op, event.Name)
					callback()
				}
			case err := <-watcher.Errors:
				errCh <- err
			}
		}
	}()

	go func() {
		for {
			select {
			case <-renameCh:
				err = waitUntilFind(file)
				if err != nil {
					log.Fatalln(err)
				}
				err = watcher.Add(file)
				if err != nil {
					log.Fatalln(err)
				}
			case <-removeCh:
				err = waitUntilFind(file)
				if err != nil {
					log.Fatalln(err)
				}
				err = watcher.Add(file)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}()

	log.Fatalln(<-errCh)
}
func waitUntilFind(filename string) error {
	for {
		time.Sleep(1 * time.Second)
		_, err := os.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				return err
			}
		}
		break
	}
	return nil
}

func MakeContextCorrelationID(parent context.Context) context.Context {
	uuID := uuid.New().String()
	corrFields := []zapcore.Field{}
	corrFields = append(corrFields, zap.Any("corr_id", uuID))
	ctx := context.WithValue(parent, logger.LOG_FIELD_KEY, corrFields)
	return ctx
}
