package report

import (
	"fmt"
	"strings"
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/golang/log_analyzer/parser"
)

type Summary struct {
	TotalLogs       int
	LevelCount      map[string]int
	FrequentMessage string
}

var (
	logSummary  = Summary{LevelCount: make(map[string]int)}
	messageFreq = make(map[string]int)
	mutex       sync.Mutex
)

func UpdateLogStats(entry parser.LogEntry) {
	logSummary.TotalLogs++
	logSummary.LevelCount[entry.Level]++
	messageFreq[entry.Message]++
}

func findMostFrequentMessage() []string {
	var frequentMessages []string
	var maxCount int

	for msg, count := range messageFreq {
		if count > maxCount {
			maxCount = count
			frequentMessages = []string{msg}
		} else if count == maxCount {
			frequentMessages = append(frequentMessages, msg)
		}

	}
	return frequentMessages
}

func GenerateReport(logEntries []parser.LogEntry) {
	for _, entry := range logEntries {
		UpdateLogStats(entry)
	}

	frequentMessages := findMostFrequentMessage()
	formattedMessages := strings.Join(frequentMessages, ", ")

	fmt.Println("\n========= Log Analysis Summary =========")
	fmt.Printf("%-25s : %d\n", "Total Logs Processed", logSummary.TotalLogs)
	fmt.Printf("%-25s : %d\n", "INFO Logs", logSummary.LevelCount["info"])
	fmt.Printf("%-25s : %d\n", "WARN Logs", logSummary.LevelCount["warning"])
	fmt.Printf("%-25s : %d\n", "ERROR Logs", logSummary.LevelCount["error"])
	fmt.Printf("%-25s : %d\n", "UNKNOWN Logs", logSummary.LevelCount["unknown"])
	fmt.Printf("%-25s : %s\n", "Most Frequent Message", formattedMessages)
	fmt.Println("========================================")

}
