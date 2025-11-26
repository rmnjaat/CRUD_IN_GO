package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Logger struct {
	module_name   string
	log_file_path string
}

func New(module_name, log_file_path string) *Logger {
	return &Logger{
		module_name:   module_name,
		log_file_path: log_file_path,
	}
}

func (l *Logger) writeLog(level, message string) {
	_, file, line, ok := runtime.Caller(2)

	if !ok {
		file = "unknown"
		line = -1
	}

	file_name := filepath.Base(file)

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	logEntry := fmt.Sprintf("[%s] [%s] [%s] [%s:%d] %s \n",
		l.module_name,
		timestamp,
		level,
		file_name,
		line,
		message,
	)

	// Write to file or use DB call to write
	if l.log_file_path != "" {
		l.write_to_file(logEntry)
	}

	// print as well
	fmt.Print(logEntry)
}

func (l *Logger) write_to_file(logentry string) {

	dir := filepath.Dir(l.log_file_path)

	if dir != "" {
		os.MkdirAll(dir, 0755)
	}

	file, err := os.OpenFile(l.log_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(logentry)
	if err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}

}

func (l Logger) Info(message string) {
	l.writeLog("INFO", message)
}

func (l Logger) Debug(message string) {
	l.writeLog("DEBUG", message)
}

func (l Logger) Warning(message string) {
	l.writeLog("WARNING", message)
}

func (l Logger) Error(message string) {

	l.writeLog("ERROR", message)
}

func (l Logger) Fatal(message string) {
	l.writeLog("FATAL", message)
}
