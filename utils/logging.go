package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Logger *log.Logger

func init() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("logs/cannoli_%s.log", timestamp)

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}

	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	Logger.Printf("=== Cannoli OS Started at %s ===", time.Now().Format(time.RFC3339))
}
