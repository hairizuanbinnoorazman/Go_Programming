package ga

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

type Client struct {
	Service    *analyticsdata.Service
	PropertyID string
}

func NewClient(ctx context.Context) (*Client, error) {
	propertyID := os.Getenv("GA_PROPERTY_ID")
	if propertyID == "" {
		return nil, fmt.Errorf("GA_PROPERTY_ID environment variable is not set")
	}

	service, err := analyticsdata.NewService(ctx, option.WithScopes(analyticsdata.AnalyticsReadonlyScope))
	if err != nil {
		return nil, fmt.Errorf("failed to create analytics service: %w", err)
	}

	return &Client{
		Service:    service,
		PropertyID: propertyID,
	}, nil
}

func (c *Client) RunReport(req *analyticsdata.RunReportRequest) (*analyticsdata.RunReportResponse, error) {
	return c.Service.Properties.RunReport(c.PropertyID, req).Do()
}
