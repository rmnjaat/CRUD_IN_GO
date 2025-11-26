package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

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
