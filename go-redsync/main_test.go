package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mandarinkb/go-redsync/config"
	"github.com/mandarinkb/go-redsync/external/redsync"
	"github.com/mandarinkb/go-redsync/logger"
)

func TestXxx(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)

	cfg := config.LoadConfig("config", "config")
	mainLog := logger.InitialLogger()

	redsync.NewClient(cfg.Redis)
	keyMutex := fmt.Sprintf("%s_mutex", cfg.RedisOption.MyOption.KeyFormat)
	mutex := redsync.NewMutex(keyMutex)

	go func() {
		tn := time.Now().String()
		fmt.Printf("%s : %v", tn, "func_1 success\n")
		PrintLog(mainLog, mutex)
		wg.Done()
	}()

	go func() {
		tn := time.Now().String()
		fmt.Printf("%s : %v", tn, "func_2 success\n")
		PrintLog(mainLog, mutex)
		wg.Done()
	}()

	go func() {
		tn := time.Now().String()
		fmt.Printf("%s : %v", tn, "func_3 success\n")
		PrintLog(mainLog, mutex)
		wg.Done()
	}()
	wg.Wait()
}
