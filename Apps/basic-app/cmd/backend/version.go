package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the version number of %s", serviceName),
		Long:  fmt.Sprintf("Print the version number of %s", serviceName),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v\n", version)
		},
	}
)
