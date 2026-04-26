package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"google.golang.org/api/analyticsdata/v1beta"
)

var (
	startDate string
	endDate   string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Google Analytics data",
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().StringVar(&startDate, "start-date", "30daysAgo", "Start date for the report (e.g. 2023-01-01, 30daysAgo, today)")
	searchCmd.PersistentFlags().StringVar(&endDate, "end-date", "today", "End date for the report (e.g. 2023-01-31, today)")
}

func renderTable(header []string, resp *analyticsdata.RunReportResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, strings.Join(header, "\t"))

	for _, row := range resp.Rows {
		var rowData []string
		for _, dv := range row.DimensionValues {
			rowData = append(rowData, dv.Value)
		}
		for _, mv := range row.MetricValues {
			rowData = append(rowData, mv.Value)
		}
		fmt.Fprintln(w, strings.Join(rowData, "\t"))
	}
	w.Flush()
}
