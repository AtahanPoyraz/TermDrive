package logger

import (
	"fmt"
	"os"
	"time"
)

// Logger writes log messages to a file with a timestamp.
// It creates a "log" directory if it does not exist and writes logs
// to a file named with the current date (YYYY-MM-DD.log).
func Logger(text string) (err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, text)

	if err := os.MkdirAll("./log", os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(fmt.Sprintf("./log/%s.log", time.Now().Format("2006-01-02")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(logMessage); err != nil {
		return err
	}

	return nil
}
