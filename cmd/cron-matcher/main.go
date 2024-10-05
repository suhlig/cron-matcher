package main

import (
	"fmt"
	"os"
	"time"

	cron "github.com/suhlig/cron-matcher"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: program <cron expression>")
		os.Exit(1)
	}

	expression := os.Args[1]
	now := time.Now()

	matches, err := cron.Matches(expression, now)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if matches {
		fmt.Println(now.Format(time.RFC3339))
		os.Exit(0)
	}

	os.Exit(1)
}
