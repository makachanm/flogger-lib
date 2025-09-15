package flogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Logger manages the connection to the flogger server.
type Logger struct {
	conn     net.Conn
	clientID int32
	mu       sync.Mutex
}

// New connects to the flogger server and returns a new Logger.
func New() (*Logger, error) {
	conn, err := net.Dial("unix", "/tmp/flogger.sock")
	if err != nil {
		return nil, fmt.Errorf("flogger: failed to connect to server: %w", err)
	}

	return &Logger{
		conn:     conn,
		clientID: int32(os.Getpid()),
	}, nil
}

// Close closes the connection to the server.
func (l *Logger) Close() error {
	if l.conn != nil {
		return l.conn.Close()
	}
	return nil
}

// send formats and writes the log message to the server.
func (l *Logger) send(msg []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var buf bytes.Buffer

	// Write client ID (int32)
	if err := binary.Write(&buf, binary.BigEndian, l.clientID); err != nil {
		return fmt.Errorf("writing client id: %w", err)
	}

	// Write message length (uint32)
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(msg))); err != nil {
		return fmt.Errorf("writing message length: %w", err)
	}

	// Write message
	if _, err := buf.Write(msg); err != nil {
		return fmt.Errorf("writing message: %w", err)
	}

	// Send data over the connection
	if _, err := l.conn.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("sending data: %w", err)
	}

	return nil
}

// Printf sends a formatted log message to the server.
func (l *Logger) Printf(format string, v ...interface{}) {
	if err := l.send([]byte(fmt.Sprintf(format, v...))); err != nil {
		fmt.Fprintf(os.Stderr, "flogger error: %v\n", err)
	}
}

// Println sends a log message to the server.
func (l *Logger) Println(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	if err := l.send([]byte(msg[:len(msg)-1])); err != nil {
		fmt.Fprintf(os.Stderr, "flogger error: %v\n", err)
	}
}

// Print sends a log message to the server.
func (l *Logger) Print(v ...interface{}) {
	if err := l.send([]byte(fmt.Sprint(v...))); err != nil {
		fmt.Fprintf(os.Stderr, "flogger error: %v\n", err)
	}
}

// --- Default Logger ---

var defaultLogger *Logger

func init() {
	var err error
	defaultLogger, err = New()
	if err != nil {
		// Silently fail. Calls to the default logger will do nothing.
		defaultLogger = nil
		return // Don't set up closer if connection failed.
	}

	// Automatically close the default logger's connection on program termination.
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan // Block until a signal is received
		CloseDefault()
	}()
}

// Printf sends a formatted log message using the default logger.
func Printf(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Printf(format, v...)
	}
}

// Println sends a log message using the default logger.
func Println(v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Println(v...)
	}
}

// Print sends a log message using the default logger.
func Print(v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Print(v...)
	}
}

// CloseDefault closes the connection for the default logger.
func CloseDefault() {
	if defaultLogger != nil {
		defaultLogger.Close()
	}
}
