package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download sign-in data from Azure",
	Run: func(cmd *cobra.Command, args []string) {
		outputPath, err := cmd.Flags().GetString("output-path")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Downloading sign-in data to %s\n", outputPath)

	},
}
