package main

import (
	"fmt"
	"time"

	flogger "github.com/makachanm/flogger-lib"
)

func main() {
	fmt.Println("--- Using Default Logger ---")
	// These will use the default logger initialized by the library.
	// Make sure the flogger-server is running.
	flogger.Println("This is a test message from the default logger.")
	flogger.Printf("This is a formatted message with a number: %d", 123)
	flogger.Print("This is a simple print message.")
	// Close the default logger's connection when your application shuts down.
	defer flogger.CloseDefault()

	// Give the server a moment to process messages before exiting
	time.Sleep(100 * time.Millisecond)
}
