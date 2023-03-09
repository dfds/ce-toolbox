package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.dfds.cloud/u/emcla/aws-sso-usage-report-gen/internal/model"
	"log"
	"os"
)

var generalCmd = &cobra.Command{
	Use:   "general",
	Short: "Get general stats",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatal(err)
		}

		stats := model.NewStats()
		data := model.LoadInputData(path)
		stats.CalcGeneralStats(data)
		fmt.Print(stats.OutputTotalAsCsv())
	},
}

var fullReportCmd = &cobra.Command{
	Use:   "full-report",
	Short: "Generate full report with all stats",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatal(err)
		}
		outputPath, err := cmd.Flags().GetString("output-path")
		if err != nil {
			log.Fatal(err)
		}

		stats := model.NewStats()
		data := model.LoadInputData(path)
		stats.CalcGeneralStats(data)
		stats.CalcCountries(data)

		if _, err = os.Stat(outputPath); os.IsNotExist(err) {
			err = os.MkdirAll(outputPath, 0770)
			if err != nil {
				log.Fatal(err)
			}
		}

		totalData := stats.OutputTotalAsCsv()
		countryData := stats.OutputUserCountryAsCsv()
		userActivityData := stats.OutputUserActivityAsCsv()

		err = os.WriteFile(fmt.Sprintf("%s/total.csv", outputPath), []byte(totalData), 0770)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(fmt.Sprintf("%s/country.csv", outputPath), []byte(countryData), 0770)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(fmt.Sprintf("%s/user.csv", outputPath), []byte(userActivityData), 0770)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func initGeneral() {
	generalCmd.Flags().StringP("input", "i", "", "Path to input file")
	generalCmd.MarkFlagRequired("input")

	fullReportCmd.Flags().StringP("input", "i", "", "Path to input file")
	fullReportCmd.Flags().StringP("output-path", "w", "output", "Path to output report files")
	fullReportCmd.MarkFlagRequired("input")
}
