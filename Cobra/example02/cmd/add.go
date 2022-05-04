package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add subcommand add all passed args.",
	Run: func(cmd *cobra.Command, args []string) {
		values := ConvertArgsToFloat64Slice(args, ErrorHandling(parseHandling))
		result := calc(values, ADD)
		fmt.Printf("%s = %.2f\n", strings.Join(args, "+"), result)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
