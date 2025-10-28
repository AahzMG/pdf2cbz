package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "PDF2CBZ",
	Short: "PDF2CBZ is a tool to convert comics in PDF format to CBZ format.",
	Long:  `PDF2CBZ is a tool to convert comics in PDF format to CBZ format.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: '%s'\n", err)
		os.Exit(1)
	}
}
