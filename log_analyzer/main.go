package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/golang/log_analyzer/parser"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/golang/log_analyzer/report"
)

func main() {
	//if flag of logdir is not provided then by default it will take ./logs directory to analyze log files
	logDir := flag.String("logdir", "./logs", "Path to the log directory containing log files")

	flag.Parse()

	if _, err := os.Stat(*logDir); os.IsNotExist(err) {
		log.Fatalf("Error: Log directory '%s' does not exist.\n", *logDir)
	}

	files, err := filepath.Glob(filepath.Join(*logDir, "*.json"))
	if err != nil {
		log.Fatalf("Error reading log files: %s", err)
	}

	if len(files) == 0 {
		fmt.Println("No log files found in the directory.")
		return
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var allLogs []parser.LogEntry

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			logEntries, err := parser.ProcessLogFile(file)
			if err != nil {
				log.Printf("Skipping file %s due to error: %v", file, err)
				return
			}

			mutex.Lock()
			allLogs = append(allLogs, logEntries...)
			mutex.Unlock()
		}(file)
	}
	wg.Wait()

	report.GenerateReport(allLogs)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Processing Time: %d ms %d µs\n", elapsedTime.Milliseconds(), elapsedTime.Microseconds()%1000)

}
