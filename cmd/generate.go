package cmd

import (
	"github.com/epentland/gohandlr/pkg/codegen"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go code",
	Long:  `Generate Go code that uses the library package.`,
	Run: func(cmd *cobra.Command, args []string) {
		codegen.GenerateCode(openapiPath)
	},
}
