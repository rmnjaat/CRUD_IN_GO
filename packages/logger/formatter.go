package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

func (l *Logger) writeLog(level, message string) {
	_, file, line, ok := runtime.Caller(2)

	if !ok {
		file = "unknown"
		line = -1
	}

	file_name := filepath.Base(file)

	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")

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
