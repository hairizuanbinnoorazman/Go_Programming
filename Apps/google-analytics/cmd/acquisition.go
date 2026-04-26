package cmd

import (
	"log"

	"google-analytics-cli/pkg/ga"

	"github.com/spf13/cobra"
	"google.golang.org/api/analyticsdata/v1beta"
)

var acquisitionCmd = &cobra.Command{
	Use:   "acquisition",
	Short: "Search for acquisition data",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ga.NewClient(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		req := &analyticsdata.RunReportRequest{
			Dimensions: []*analyticsdata.Dimension{{Name: "sessionSourceMedium"}},
			Metrics: []*analyticsdata.Metric{
				{Name: "totalUsers"},
				{Name: "newUsers"},
			},
			DateRanges: []*analyticsdata.DateRange{{StartDate: startDate, EndDate: endDate}},
		}

		resp, err := client.RunReport(req)
		if err != nil {
			log.Fatal(err)
		}

		renderTable([]string{"Source / Medium", "Total Users", "New Users"}, resp)
	},
}

func init() {
	searchCmd.AddCommand(acquisitionCmd)
}
