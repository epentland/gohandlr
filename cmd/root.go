package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var openapiPath string

var rootCmd = &cobra.Command{
	Use:   "gohandlr",
	Short: "A CLI tool for code generation",
	Long:  `A CLI tool for generating Go code that uses a library package.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands here
	rootCmd.AddCommand(generateCmd)

	// Define flags
	generateCmd.Flags().StringVarP(&openapiPath, "openapi", "o", "", "Path to the openapi.yaml file (required)")

	// Mark flags as required
	generateCmd.MarkFlagRequired("openapi")
}
