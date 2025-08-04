// Package core provides the core functionality for reading logs from files.
package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// LogReader defines the contract for reading logs.
type LogReader interface {
	// Tail reads the log file in real-time and sends the new lines to the
	// provided VWriter or prints them to the terminal. The operation can be interrupted
	// by sending a signal through the stopChan channel.
	Tail(filePath string, stopChan <-chan struct{}) error
}

// FileLogReader implements the LogReader interface by reading from a file.
type FileLogReader struct {
	// pollInterval is the polling interval to check for new lines.
	pollInterval time.Duration
}

// NewFileLogReader creates a new instance of FileLogReader.
// The polling interval is read from the LOGZ_TAIL_POLL_INTERVAL environment variable (in milliseconds),
// or defaults to 500ms.
func NewFileLogReader() *FileLogReader {
	intervalMs := 500 // default in milliseconds
	if val := os.Getenv("LOGZ_TAIL_POLL_INTERVAL"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			intervalMs = parsed
		} else {
			log.Printf("Invalid LOGZ_TAIL_POLL_INTERVAL value, using default 500ms: %v", err)
		}
	}
	return &FileLogReader{
		pollInterval: time.Duration(intervalMs) * time.Millisecond,
	}
}

// Tail follows the log file from the end and prints new lines as they are added.
// The stopChan channel allows interrupting the operation (e.g., via Ctrl+C).
func (fr *FileLogReader) Tail(filePath string, stopChan <-chan struct{}) error {
	// Open the log file
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	// Position the pointer at the end of the file to read only new lines
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek to the end of the file: %w", err)
	}

	reader := bufio.NewReader(f)

	// Main loop to read new lines
	for {
		select {
		case <-stopChan:
			log.Println("Tail operation interrupted by stop signal")
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(fr.pollInterval)
					continue
				}
				return fmt.Errorf("error reading log file: %w", err)
			}
			// Print the line immediately; can be adapted to send to another channel if needed.
			fmt.Print(line)
		}
	}
}
