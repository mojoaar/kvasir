package cmd

import (
	"github.com/spf13/cobra"
	"kvasir-cli/cmd/create"
	"kvasir-cli/cmd/export"
	"kvasir-cli/cmd/list"
	"kvasir-cli/cmd/search"
	"kvasir-cli/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "kvasir",
	Short: "Kvasir — the kernel of your mind",
	Long:  "A CLI for interacting with the Kvasir knowledge base.",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.ServerURL, "server", "http://localhost:8080", "Kvasir server URL")
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(search.Cmd)
	rootCmd.AddCommand(export.Cmd)
	rootCmd.AddCommand(create.Cmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
