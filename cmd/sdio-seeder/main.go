package main

import (
	"os"
	"sync"
	"time"

	"github.com/OverlyDev/go-sdio-seeder/internal/logger"
	"github.com/OverlyDev/go-sdio-seeder/internal/sdio"
	"github.com/OverlyDev/go-sdio-seeder/internal/torrent"
)

const updateInterval = 15 * time.Minute
const statusInterval = 1 * time.Minute

var helper sdio.SdioHelper

func init() {
	// logger.EnableDebugLogging()
	helper = sdio.NewHelper()
	torrent.Setup()
}

func updateLoop() {
	for {
		time.Sleep(updateInterval)
		helper.Update()
	}
}

func statusLoop() {
	for {
		time.Sleep(statusInterval)
		helper.Status()
	}
}

func main() {
	if err := helper.Start(); err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go updateLoop()
	go statusLoop()

	wg.Wait()
}
