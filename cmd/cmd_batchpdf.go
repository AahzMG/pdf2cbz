package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var batchCmd = &cobra.Command{
	Use:     "batch",
	Aliases: []string{"c"},
	Short:   "Converts all PDFs in folder enter path",
	Long:    "Converts all PDFs in folder enter path 'pdf2cbz batch PathToFolder'",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Conversion of %s : %s.\n\n", args[0], batchPDF2CBZ(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)
}
