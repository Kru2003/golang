# Log Analyzer

## Overview
Log Analyzer is a Go-based tool designed to process and analyze log files from a distributed system. It can parse multiple JSON log files simultaneously, extract insights, and generate summary reports.

## Features
- Reads and processes multiple log files concurrently.
- Supports structured log entries with dynamic fields.
- Aggregates log levels (INFO, WARN, ERROR).
- Generates a report summarizing key insights.

## Project Structure
log-analyzer/
│── go.mod
│── main.go
│── parser/
│ ├── parser.go
│── report/
│ ├── report.go
│── logs/
│ ├── auth.json
│ ├── database.json
│ ├── file.json
│── README.md


## Installation & Setup

1. Clone the repository:
  git clone git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/golang/log_analyzer
  cd log-analyzer

2. Initialize Go modules (if not already initialized):
  go mod tidy


## Run the log analyzer with default log files
go run main.go

## Run the log analyzer with specific log files
go run main.go -logdir=/path/to/logdir/
