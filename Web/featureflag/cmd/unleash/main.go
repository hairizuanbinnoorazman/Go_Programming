package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	unleash "github.com/Unleash/unleash-client-go/v3"
	unleashProvider "github.com/open-feature/go-sdk-contrib/providers/unleash/pkg"
	"github.com/open-feature/go-sdk/pkg/openfeature"
)

var (
	tpl = `<!DOCTYPE html>
<head>
<style>
body {
	background-color: {{ .BgColor }};
}

h1 {
	color: {{ .TitleColor }};
}
</style>
</head>
<html>
	<body>
		<h1>{{ .MainTitle }}</h1>
		{{ if ne .User "" }}
		<p>Pretending to be the following user: {{ .User }}</p>
		{{ end }}
	</body>
</html>`
)

func main() {
	apiKey, exists := os.LookupEnv("API_KEY")
	if !exists {
		log.Fatal("API_KEY environment variable not set")
	}
	// Example: "http://localhost:4242/api/"
	ffPlatform, exists := os.LookupEnv("FF_URL")
	if !exists {
		ffPlatform = "http://localhost:4242/api/"
	}

	providerConfig := unleashProvider.ProviderConfig{
		Options: []unleash.ConfigOption{
			// unleash.WithListener(&unleash.DebugListener{}),
			unleash.WithAppName("my-application"),
			unleash.WithRefreshInterval(5 * time.Second),
			unleash.WithMetricsInterval(5 * time.Second),
			unleash.WithUrl(ffPlatform),
			unleash.WithCustomHeaders(http.Header{"Authorization": {apiKey}}),
		},
	}
	provider, err := unleashProvider.NewProvider(providerConfig)
	if err != nil {
		fmt.Printf("Error: unable to create provider %v", err)
		os.Exit(1)
	}
	err = provider.Init(openfeature.EvaluationContext{})
	if err != nil {
		fmt.Printf("Error: unable to create provider %v", err)
		os.Exit(1)
	}

	openfeature.SetProvider(provider)
	ofClient := openfeature.NewClient("my-app")

	log.Println("begin server")

	http.Handle("/", &index{OpenFeatureCl: ofClient})
	s := &http.Server{
		Addr: "0.0.0.0:8080",
	}
	log.Fatal(s.ListenAndServe())
}

type index struct {
	OpenFeatureCl *openfeature.Client
}

func (i *index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start index handler")
	defer log.Println("end index handler")

	pretendUser := r.URL.Query().Get("pretend")
	if pretendUser != "" {
		log.Printf("pretending to be the following user: %v", pretendUser)
	}

	evalCtx := openfeature.NewEvaluationContext(
		"",
		map[string]interface{}{
			"UserId": pretendUser,
		},
	)
	webpageTitle, err := i.OpenFeatureCl.StringValue(context.Background(), "webpage-title", "Welcome", evalCtx)
	if err != nil {
		log.Printf("error occured - but code can continue fine: %v\n", err)
	}
	backgroundColor, err := i.OpenFeatureCl.StringValue(context.Background(), "background-color", "transparent", evalCtx)
	if err != nil {
		log.Printf("error occured - but code can continue fine: %v\n", err)
	}
	titleColor, err := i.OpenFeatureCl.StringValue(context.Background(), "title-color", "black", evalCtx)
	if err != nil {
		log.Printf("error occured - but code can continue fine: %v\n", err)
	}

	t, err := template.New("(webpage").Parse(tpl)
	if err != nil {
		log.Printf("error occured when parsing: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse"))
		return
	}

	data := struct {
		MainTitle  string
		BgColor    string
		User       string
		TitleColor string
	}{
		MainTitle:  webpageTitle,
		BgColor:    backgroundColor,
		User:       pretendUser,
		TitleColor: titleColor,
	}

	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("error occured when parsing: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse"))
		return
	}
}
