package ga

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

var (
	clientID     = os.Getenv("GA_CLIENT_ID")
	clientSecret = os.Getenv("GA_CLIENT_SECRET")
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
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GA_CLIENT_ID and GA_CLIENT_SECRET environment variables must be set")
	}

	tokenSource, err := getTokenSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	service, err := analyticsdata.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("failed to create analytics service: %w", err)
	}

	return &Client{
		Service:    service,
		PropertyID: propertyID,
	}, nil
}

func (c *Client) RunReport(req *analyticsdata.RunReportRequest) (*analyticsdata.RunReportResponse, error) {
	return c.Service.Properties.RunReport("properties/"+c.PropertyID, req).Do()
}

func oauthConfig(redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       []string{analyticsdata.AnalyticsReadonlyScope},
	}
}

func tokenCachePath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "ga-cli", "token.json")
}

func loadCachedToken() (*oauth2.Token, error) {
	data, err := os.ReadFile(tokenCachePath())
	if err != nil {
		return nil, err
	}
	var tok oauth2.Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func saveCachedToken(tok *oauth2.Token) error {
	p := tokenCachePath()
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	data, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func getTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	tok, err := loadCachedToken()
	if err == nil && tok.Valid() {
		cfg := oauthConfig("http://localhost:0")
		return cfg.TokenSource(ctx, tok), nil
	}
	if err == nil && tok.RefreshToken != "" {
		cfg := oauthConfig("http://localhost:0")
		ts := cfg.TokenSource(ctx, tok)
		newTok, err := ts.Token()
		if err == nil {
			_ = saveCachedToken(newTok)
			return cfg.TokenSource(ctx, newTok), nil
		}
	}

	tok, err = doOAuthFlow(ctx)
	if err != nil {
		return nil, err
	}
	_ = saveCachedToken(tok)
	cfg := oauthConfig("http://localhost:0")
	return cfg.TokenSource(ctx, tok), nil
}

func doOAuthFlow(ctx context.Context) (*oauth2.Token, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start local server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURL := fmt.Sprintf("http://localhost:%d", port)

	cfg := oauthConfig(redirectURL)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			fmt.Fprintf(w, "Error: no authorization code received.")
			return
		}
		codeCh <- code
		fmt.Fprintf(w, "Authorization successful! You can close this window.")
	})
	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(listener) }()
	defer srv.Shutdown(ctx)

	authURL := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Open this URL in your browser to authorize:\n\n%s\n\nWaiting for authorization...\n", authURL)

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}
	return tok, nil
}
