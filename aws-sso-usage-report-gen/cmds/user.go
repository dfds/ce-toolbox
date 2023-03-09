package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.dfds.cloud/u/emcla/aws-sso-usage-report-gen/internal/model"
	"log"
)

var userActivityCmd = &cobra.Command{
	Use:   "user-activity",
	Short: "Total user sign-ins by user",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatal(err)
		}

		stats := model.NewStats()
		data := model.LoadInputData(path)
		stats.CalcGeneralStats(data)
		fmt.Print(stats.OutputUserActivityAsCsv())
	},
}

var userUniqueCountryCmd = &cobra.Command{
	Use:   "user-country",
	Short: "Country stats",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatal(err)
		}

		stats := model.NewStats()
		data := model.LoadInputData(path)
		stats.CalcCountries(data)
		fmt.Print(stats.OutputUserCountryAsCsv())
	},
}

func initUser() {
	userActivityCmd.Flags().StringP("input", "i", "", "Path to input file")
	userActivityCmd.MarkFlagRequired("input")
	userUniqueCountryCmd.Flags().StringP("input", "i", "", "Path to input file")
	userUniqueCountryCmd.MarkFlagRequired("input")
}
