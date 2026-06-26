package export

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"kvasir-cli/pkg/api"
	"kvasir-cli/pkg/config"
)

var outputFile string

var Cmd = &cobra.Command{
	Use:   "export [id]",
	Short: "Export a note as markdown",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid note ID: %s", args[0])
		}

		client := api.NewClient(config.ServerURL)
		note, err := client.GetNote(id)
		if err != nil {
			return fmt.Errorf("export failed: %w", err)
		}

		filename := outputFile
		if filename == "" {
			filename = fmt.Sprintf("%s.md", sanitizeFilename(note.Title))
		}

		if err := os.WriteFile(filename, []byte(note.Content), 0644); err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		fmt.Printf("Exported note %d to %s\n", note.ID, filename)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: <title>.md)")
}

func sanitizeFilename(name string) string {
	result := make([]byte, 0, len(name))
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			result = append(result, c)
		} else if c == ' ' {
			result = append(result, '-')
		}
	}
	if len(result) == 0 {
		return "note"
	}
	return string(result)
}
