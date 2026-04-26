package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ga-cli",
	Short: "Google Analytics 4 CLI tool",
	Long:  `A command line interface for querying Google Analytics 4 data using ADC.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags can be added here if needed
}
