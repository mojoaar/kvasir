package create

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"kvasir-cli/pkg/api"
	"kvasir-cli/pkg/config"
)

var noteContent string

var Cmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		content := noteContent

		if content == "" {
			content = fmt.Sprintf("# %s\n\n", title)
		}

		client := api.NewClient(config.ServerURL)
		note, err := client.CreateNote(title, content)
		if err != nil {
			return fmt.Errorf("create failed: %w", err)
		}

		fmt.Printf("Created note %d: %s\n", note.ID, note.Title)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&noteContent, "content", "c", "", "Note content in markdown (default: # <title>)")
}
