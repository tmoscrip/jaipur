package logger

import (
	"fmt"
	"os"
	"time"
)

// Error writes error messages to a log file with a timestamp.
func Error(err error) {
	if err == nil {
		return
	}

	// Open the log file in append mode, create it if it doesn't exist
	file, fileErr := os.OpenFile("messages.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		// If opening the file fails, we can panic or handle it as necessary
		// However, avoid logging the error to the console if you're in a text UI
		panic(fmt.Sprintf("failed to open log file: %v", fileErr))
	}
	defer file.Close()

	// Prepare the log message with a timestamp
	logMessage := fmt.Sprintf("%s: %v\n", time.Now().Format(time.RFC3339), err)

	// Write the log message to the file
	if _, writeErr := file.WriteString(logMessage); writeErr != nil {
		// If writing to the file fails, handle the error as necessary
		panic(fmt.Sprintf("failed to write to log file: %v", writeErr))
	}
}

// Message writes a message string to a log file with a timestamp.
func Message(message string) {
	// Open the log file in append mode, create it if it doesn't exist
	file, fileErr := os.OpenFile("messages.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		// If opening the file fails, we can panic or handle it as necessary
		// However, avoid logging the error to the console if you're in a text UI
		panic(fmt.Sprintf("failed to open log file: %v", fileErr))
	}
	defer file.Close()
	// Prepare the log message with a timestamp
	logMessage := fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), message)
	// Write the log message to the file
	if _, writeErr := file.WriteString(logMessage); writeErr != nil {
		// If writing to the file fails, handle the error as necessary
		panic(fmt.Sprintf("failed to write to log file: %v", writeErr))
	}
}
