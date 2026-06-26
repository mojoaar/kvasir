package search

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"kvasir-cli/pkg/api"
	"kvasir-cli/pkg/config"
)

var limit int

var Cmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Full-text search across notes",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		client := api.NewClient(config.ServerURL)

		results, err := client.Search(query, limit)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("No results found.")
			return nil
		}

		for _, r := range results {
			fmt.Printf("%s (%s) — ID: %d\n", r.Title, noteType(r.IsFolder), r.ID)
			if r.Snippet != "" {
				snippet := strings.ReplaceAll(r.Snippet, "<mark>", "\033[33m")
				snippet = strings.ReplaceAll(snippet, "</mark>", "\033[0m")
				fmt.Printf("  %s\n", snippet)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	Cmd.Flags().IntVar(&limit, "limit", 10, "Max results to return")
}

func noteType(isFolder bool) string {
	if isFolder {
		return "folder"
	}
	return "note"
}
