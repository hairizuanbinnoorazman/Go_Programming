package cmd

import (
	"log"

	"google-analytics-cli/pkg/ga"

	"github.com/spf13/cobra"
	"google.golang.org/api/analyticsdata/v1beta"
)

var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Search for pages and screens",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ga.NewClient(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		req := &analyticsdata.RunReportRequest{
			Dimensions: []*analyticsdata.Dimension{{Name: "pagePath"}},
			Metrics:    []*analyticsdata.Metric{{Name: "screenPageViews"}},
			DateRanges: []*analyticsdata.DateRange{{StartDate: startDate, EndDate: endDate}},
		}

		resp, err := client.RunReport(req)
		if err != nil {
			log.Fatal(err)
		}

		renderTable([]string{"Page Path", "Views"}, resp)
	},
}

func init() {
	searchCmd.AddCommand(pagesCmd)
}
