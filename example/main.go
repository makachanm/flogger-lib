package main

import (
	"fmt"
	"time"

	flogger "github.com/makachanm/flogger-lib"
)

func main() {
	fmt.Println("--- Using Default Logger ---")
	flogger.Println("This is a test message from the default logger.")
	flogger.Printf("This is a formatted message with a number: %d", 123)
	flogger.Print("This is a simple print message.")

	time.Sleep(100 * time.Millisecond)
}
