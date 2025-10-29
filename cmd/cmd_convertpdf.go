package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"c"},
	Short:   "Converts PDF enter path to PDF",
	Long:    "Converts PDF enter path to PDF 'pdf2cbz convert myfile.pdf'",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Conversion of %s : %s.\n\n", args[0], convertPDF2CBZ(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
