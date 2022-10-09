package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func main() {
	log.Println("Begin server")
	bucketName := os.Getenv("BUCKET_NAME")
	log.Printf("Bucket name provided: %v\n", bucketName)

	cl, err := storage.NewClient(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to start server :: %v\n", err))
	}

	r := mux.NewRouter()
	r.Handle("/create-report", &CreateReport{
		BucketName: bucketName,
		StorageSvc: cl,
	}).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}

type CreateReport struct {
	BucketName string
	StorageSvc *storage.Client
}

type createReportRequest struct {
	Title            string `json:"title"`
	TemplateFileName string `json:"template_file_name"`
	Image            string `json:"image"`
	Description      string `json:"description"`
}

func (c *CreateReport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("createReport request started")
	defer log.Println("createReport request ended")

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to read body request :: %v", err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var createReportReq createReportRequest
	json.Unmarshal(raw, &createReportReq)
	outputFileName := strings.ReplaceAll(createReportReq.Title, " ", "_") + ".pdf"

	// Download images
	err = c.downloadFile(ctx, createReportReq.Image)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unexpected error :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(createReportReq.Image)

	// Download template
	c.downloadFile(ctx, createReportReq.TemplateFileName)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unexpected error :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(createReportReq.TemplateFileName)

	md := goldmark.New(
		goldmark.WithRenderer(pdf.New(
			pdf.WithImageFS(os.DirFS(".")),
		)),
	)

	f, _ := os.Create(outputFileName)
	defer os.Remove(outputFileName)

	// source, _ := ioutil.ReadFile(createReportReq.TemplateFileName)
	tmpl, err := template.New(createReportReq.TemplateFileName).ParseFiles(createReportReq.TemplateFileName)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unexpected error :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	err = tmpl.Execute(foo, createReportReq)
	foo.Flush()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unexpected error :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	doc := md.Parser().Parse(text.NewReader(b.Bytes()))
	err = md.Renderer().Render(f, b.Bytes(), doc)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to render report :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.uploadFile(ctx, outputFileName)

	type createReportResponse struct {
		GeneratedReportName string `json:"generated_report_name"`
	}

	createReportRespRaw, _ := json.Marshal(createReportResponse{GeneratedReportName: outputFileName})
	w.Write(createReportRespRaw)
	w.WriteHeader(http.StatusCreated)
}

func (c *CreateReport) downloadFile(ctx context.Context, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	o := c.StorageSvc.Bucket(c.BucketName).Object(fileName)
	rc, err := o.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to create storage reader")
	}
	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := rc.Close(); err != nil {
		return fmt.Errorf("reader.Close: %v", err)
	}
	return nil
}

func (c *CreateReport) uploadFile(ctx context.Context, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	o := c.StorageSvc.Bucket(c.BucketName).Object(fileName)
	wc := o.NewWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to create storage reader")
	}
	if _, err := io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("reader.Close: %v", err)
	}
	return nil
}
