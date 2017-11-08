/*
There are 3 main pages

1. Main page - index page. This will contain a link to direct user to login page
2. Login page - Redirect user to server
3. Callback page -> Get code -> Exchange for token -> Redirect
4. Page for registered users
 */

package main

import (
	"log"
	"net/http"
	"io"
	"golang.org/x/oauth2"
	"google.golang.org/api/analyticsreporting/v4"
	"context"
)

var conf = &oauth2.Config{
	ClientID: "--",
	ClientSecret: "--",
	Scopes: []string{"https://www.googleapis.com/auth/analytics"},
	Endpoint: oauth2.Endpoint {
		AuthURL: "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL: "https://www.googleapis.com/oauth2/v4/token",
	},
	RedirectURL: "http://localhost:3000/callback",
}


type indexHandler struct {}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Begin index handler")
	defer log.Println("End of index handler")
	io.WriteString(w, "Index Page")
}

type loginHandler struct{}

func (h loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Login Handler - Will direct to localhost/callback endpoint")
	defer log.Println("Login Handler - End of login handler.")

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, 301)
}

func Miao(ctx context.Context, token *oauth2.Token) context.Context {
	ctx = context.WithValue(ctx, int(42), token)
	return ctx
}

type callbackHandler struct{}

func (h callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	ctx := r.Context()

	log.Print(string(code))

	tok, err := conf.Exchange(ctx, code)

	log.Println(tok.AccessToken)
	log.Println(tok.RefreshToken)
	log.Println(tok.Expiry)
	log.Println(tok.TokenType)

	if err != nil {
		log.Println( err.Error())
	}

	// Try out client
	client := conf.Client(ctx, tok)
	analyticsreportingService, err := analyticsreporting.New(client)
	if err != nil {
		log.Println(err.Error())
	}
	singleDateRange := &analyticsreporting.DateRange{"2017-02-28", "2017-02-01", []string{}, []string{}}
	multiDateRange := []*analyticsreporting.DateRange{singleDateRange}

	singleMetrics := &analyticsreporting.Metric{Expression:"ga:users"}
	multipleMetrics := []*analyticsreporting.Metric{singleMetrics}

	singleReportRequest := analyticsreporting.ReportRequest{ViewId:"--", DateRanges:multiDateRange, Metrics:multipleMetrics}
	multipleReportRequests := []*analyticsreporting.ReportRequest{&singleReportRequest}

	getReportRequest := &analyticsreporting.GetReportsRequest{multipleReportRequests,[]string{}, []string{}}
	reportResponse, err := analyticsreportingService.Reports.BatchGet(getReportRequest).Do()
	if err != nil {
		log.Println(err.Error())
	}
	singleResponseReport := reportResponse.Reports[0]
	log.Println(singleResponseReport.Data.RowCount)
	for _, row := range singleResponseReport.Data.Rows {
		singleRow, _ := row.MarshalJSON()
		log.Println(string(singleRow[:]))
		for _, raw := range row.Metrics {
			log.Println(raw.Values)
		}
	}


	io.WriteString(w, "Testing application")
	//http.Redirect(w, r.WithContext(Miao(ctx, tok)), "/users", 301)
}

type userHandler struct{}

func (h userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("User Handler activated!!")
	defer log.Println("User Handler closed")

	lol := r.Context().Value(int(42))
	log.Println(lol)
	//
	//if lol == nil {
	//log.Println("tok is nil")
	//}
	////log.Println(tok.AccessToken)
	//
	////data := fmt.Sprintf("%v, %v, %v, %v", tok.AccessToken, tok.RefreshToken, tok.Expiry, tok.TokenType)
	//
	io.WriteString(w, "Testing application 2")
}

type randomHandler struct{}

func (h randomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "lol", "lollol")
	r = r.WithContext(ctx)

	http.Redirect(w, r, "/lol", http.StatusTemporaryRedirect)
}

type lolHandler struct{}

func (h lolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	haha, err := ctx.Value("lol").(string)
	if !err {
		log.Println("Could not find value")
	}
	io.WriteString(w, haha)
}


func Decorator(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ctx := r.Context()
		ctx = context.WithValue(ctx, "lol", "lol2 lol2")
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Begin serving temporary application on port 3000")
	http.Handle("/", indexHandler{})
	http.Handle("/login", loginHandler{})
	http.Handle("/callback", callbackHandler{})
	http.Handle("/users", userHandler{})
	http.Handle("/random", randomHandler{}) // This call  is to show how context propagation don't work
	http.Handle("/lol", lolHandler{}) // If called for random, no value is passed through, a redirect is like a new request.
	http.Handle("/decor", Decorator(lolHandler{}))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
