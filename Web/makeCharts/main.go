package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
)

//go:embed static/*
var content embed.FS

type data struct {
	Title      string   `json:"title"`
	XaxisTitle string   `json:"x_axis_title"`
	Labels     []string `json:"labels"`
	Data       []int    `json:"data"`
}

func main() {
	bucketName := os.Getenv("BUCKET_NAME")

	r := mux.NewRouter()
	r.Handle("/screenshot", &ScreenshotIt{
		BucketName: bucketName,
	}).Methods("POST")
	r.Handle("/struct", &HandleViaStruct{}).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}

type ScreenshotIt struct {
	BucketName string
}

func (s *ScreenshotIt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start screenshot")
	defer log.Println("end screenshot")
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("bad response"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var d data
	json.Unmarshal(raw, &d)

	u, _ := url.Parse("http://localhost:8080/struct")
	values := u.Query()
	values.Set("title", d.Title)
	values.Set("x_axis_title", d.XaxisTitle)
	values.Set("labels", strings.Join(d.Labels, ","))
	var convertedData []string
	for _, v := range d.Data {
		convertedData = append(convertedData, strconv.Itoa(v))
	}
	values.Set("data", strings.Join(convertedData, ","))
	u.RawQuery = values.Encode()
	structURL := u.String()
	log.Println(structURL)

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(structURL, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	fileName := strings.ReplaceAll(d.Title+".jpeg", " ", "_")
	if err := ioutil.WriteFile(fileName, buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Println("screenshot file created")
	err = uploadFile(s.BucketName, fileName, fileName)
	if err != nil {
		w.Write([]byte("failed to store chart"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type resp struct {
		FileName string `json:"filename"`
	}

	outputResp := resp{FileName: fileName}
	outputRaw, _ := json.Marshal(outputResp)

	w.Write(outputRaw)
	w.WriteHeader(http.StatusCreated)
}

func uploadFile(bucket, object, filename string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	// acl := client.Bucket(bucket).Object(object).ACL()
	// if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
	// 	return fmt.Errorf("ACLHandle.Set: %v", err)
	// }
	return nil
}

type HandleViaStruct struct{}

func (*HandleViaStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	defer log.Print("End hello world request")
	var d data
	d.Title = r.URL.Query().Get("title")
	d.XaxisTitle = r.URL.Query().Get("x_axis_title")
	d.Labels = strings.Split(r.URL.Query().Get("labels"), ",")
	strData := strings.Split(r.URL.Query().Get("data"), ",")

	var intData []int
	for _, v := range strData {
		vv, _ := strconv.Atoi(v)
		intData = append(intData, vv)
	}
	d.Data = intData

	tmpl := template.Must(template.ParseFS(content, "static/chart.html"))
	tmpl.Execute(w, d)
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	log.Println("Start fullScreenshot function")
	defer log.Println("End fullScreenshot function")
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}
