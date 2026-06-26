package list

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List notes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing notes...")
	},
}
