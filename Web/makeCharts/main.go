package main

import (
	"context"
	"embed"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

//go:embed static/*
var content embed.FS

func main() {
	http.Handle("/screenshot", &ScreenshotIt{})
	http.Handle("/struct", &HandleViaStruct{})
	http.ListenAndServe(":8080", nil)
}

type ScreenshotIt struct{}

func (*ScreenshotIt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start screenshot")
	defer log.Println("end screenshot")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(`http://localhost:8080/struct`, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

type HandleViaStruct struct{}

func (*HandleViaStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	defer log.Print("End hello world request")
	type data struct {
		Title      string
		XaxisTitle string
		Labels     []string
		Data       []int
	}

	tmpl := template.Must(template.ParseFS(content, "static/chart.html"))
	tmpl.Execute(w, data{Title: "Sales Figures", XaxisTitle: "Months", Labels: []string{"Jan", "Feb"}, Data: []int{10, 20}})
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	log.Println("Start fullScreenshot function")
	defer log.Println("End fullScreenshot function")
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}
