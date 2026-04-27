package cmd

import (
	"log"

	"google-analytics-cli/pkg/ga"

	"github.com/spf13/cobra"
	"google.golang.org/api/analyticsdata/v1beta"
)

var pagesByCountryCmd = &cobra.Command{
	Use:   "pages-by-country",
	Short: "Search for pages by country",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ga.NewClient(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		req := &analyticsdata.RunReportRequest{
			Dimensions: []*analyticsdata.Dimension{
				{Name: "pagePath"},
				{Name: "country"},
			},
			Metrics:    []*analyticsdata.Metric{{Name: "screenPageViews"}},
			DateRanges: []*analyticsdata.DateRange{{StartDate: startDate, EndDate: endDate}},
		}

		resp, err := client.RunReport(req)
		if err != nil {
			log.Fatal(err)
		}

		renderTable([]string{"Page Path", "Country", "Views"}, resp)
	},
}

func init() {
	searchCmd.AddCommand(pagesByCountryCmd)
}
