package logger

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

func (l *Logger) Info(message string) {
	l.writeLog("INFO", message)
}

func (l *Logger) Debug(message string) {
	l.writeLog("DEBUG", message)
}

func (l *Logger) Warning(message string) {
	l.writeLog("WARNING", message)
}

func (l *Logger) Error(message string) {

	l.writeLog("ERROR", message)
}

func (l *Logger) Fatal(message string) {
	l.writeLog("FATAL", message)
}
