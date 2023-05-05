package main

import (
	"fmt"
	"time"

	rd "github.com/go-redsync/redsync/v4"
	"github.com/mandarinkb/go-redsync/config"
	"github.com/mandarinkb/go-redsync/external/redsync"
	"github.com/mandarinkb/go-redsync/logger"
)

func main() {
	cfg := config.LoadConfig("config", "config")
	mainLog := logger.InitialLogger()

	redsync.NewClient(cfg.Redis)
	keyMutex := fmt.Sprintf("%s_mutex", cfg.RedisOption.MyOption.KeyFormat)
	mutex := redsync.NewMutex(keyMutex)

	PrintLog(mainLog, mutex)
}

func PrintLog(myLog *logger.Logger, mutex *rd.Mutex) {
	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	// myLog.Info("run main success")
	tn := time.Now().String()
	fmt.Printf("%s : %v", tn, "print log success\n")

	// Release the lock so other processes or threads can obtain a lock.
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}
