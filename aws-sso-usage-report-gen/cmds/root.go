package cmds

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "asurg",
	Short: "Generate stats for usage of AWS SSO through AAD",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func initRoot() {
	// Args
	rootCmd.PersistentFlags().StringP("output", "o", "csv", "Output format. Supports: csv")

	// Commands
	rootCmd.AddCommand(fullReportCmd)
	rootCmd.AddCommand(generalCmd)
	rootCmd.AddCommand(userActivityCmd)
	rootCmd.AddCommand(userUniqueCountryCmd)
}

func RunCmd() {
	initUser()
	initGeneral()
	initRoot()
	rootCmd.Execute()
}
