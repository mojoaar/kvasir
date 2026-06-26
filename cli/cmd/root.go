package cmd

import (
	"github.com/spf13/cobra"
	"kvasir-cli/cmd/export"
	"kvasir-cli/cmd/list"
	"kvasir-cli/cmd/search"
)

var rootCmd = &cobra.Command{
	Use:   "kvasir",
	Short: "Kvasir — the kernel of your mind",
	Long:  "A CLI for interacting with the Kvasir knowledge base.",
}

func init() {
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(search.Cmd)
	rootCmd.AddCommand(export.Cmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
