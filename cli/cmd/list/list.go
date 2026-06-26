package list

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"kvasir-cli/pkg/api"
	"kvasir-cli/pkg/config"
)

var vaultID int64

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(config.ServerURL)
		notes, err := client.ListNotes()
		if err != nil {
			return fmt.Errorf("list failed: %w", err)
		}

		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tTITLE\tUPDATED")
		for _, n := range notes {
			noteType := "note"
			if n.IsFolder {
				noteType = "folder"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", n.ID, noteType, n.Title, n.UpdatedAt)
		}
		w.Flush()

		return nil
	},
}

func init() {
	Cmd.Flags().Int64Var(&vaultID, "vault", 0, "Filter by vault ID")
}
