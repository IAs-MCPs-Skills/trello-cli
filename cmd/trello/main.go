package main

import "os"

// Build metadata — set via ldflags at build time.
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(handleError(os.Stdout, err))
	}
}
