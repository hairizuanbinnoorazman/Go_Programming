package cmd

import (
	"log"

	"google-analytics-cli/pkg/ga"

	"github.com/spf13/cobra"
	"google.golang.org/api/analyticsdata/v1beta"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Search for events",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ga.NewClient(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		req := &analyticsdata.RunReportRequest{
			Dimensions: []*analyticsdata.Dimension{{Name: "eventName"}},
			Metrics:    []*analyticsdata.Metric{{Name: "eventCount"}},
			DateRanges: []*analyticsdata.DateRange{{StartDate: startDate, EndDate: endDate}},
		}

		resp, err := client.RunReport(req)
		if err != nil {
			log.Fatal(err)
		}

		renderTable([]string{"Event Name", "Count"}, resp)
	},
}

func init() {
	searchCmd.AddCommand(eventsCmd)
}
