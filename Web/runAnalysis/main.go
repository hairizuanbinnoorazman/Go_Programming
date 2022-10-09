package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
)

var ExpectedColumns = []string{"Date", "Product Name", "Price", "Qty"}

func main() {
	log.Println("Begin server")
	bucketName := os.Getenv("BUCKET_NAME")
	log.Printf("Bucket name provided: %v\n", bucketName)

	cl, err := storage.NewClient(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to start server :: %v\n", err))
	}

	r := mux.NewRouter()
	r.Handle("/run-analysis", &RunAnalysis{
		BucketName: bucketName,
		StorageSvc: cl,
	}).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}

type RunAnalysis struct {
	BucketName string
	StorageSvc *storage.Client
}

func (ra *RunAnalysis) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("runAnalysis request started")
	defer log.Println("runAnalysis request ended")

	type runAnalysisRequest struct {
		SourceData string `json:"source_data"`
	}
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to read body request :: %v", err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var runAnalysisReq runAnalysisRequest
	json.Unmarshal(raw, &runAnalysisReq)

	err = ra.downloadFile(ctx, runAnalysisReq.SourceData)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unexpected error :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(runAnalysisReq.SourceData)

	f, err := os.Open(runAnalysisReq.SourceData)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unexpected error :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unexpected error :: %v\n", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Need to optimize
	err = ra.checkValidQuick(records, ExpectedColumns)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("dataset error :: %v\n", err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summed, err := ra.productSummer(records)

	p := make(PairList, len(summed))
	i := 0
	for k, v := range summed {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	limitedP := p
	if len(summed) > 10 {
		limitedP = limitedP[:10]
	}

	type runAnalysisResponse struct {
		Products []string `json:"products"`
		Revenue  []int    `json:"revenue"`
	}

	var runAnalysisResp runAnalysisResponse

	for _, v := range limitedP {
		runAnalysisResp.Products = append(runAnalysisResp.Products, v.Key)
		runAnalysisResp.Revenue = append(runAnalysisResp.Revenue, v.Value)
	}

	rawResp, _ := json.Marshal(runAnalysisResp)
	w.Write(rawResp)
	w.WriteHeader(http.StatusOK)
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func (r *RunAnalysis) productSummer(records [][]string) (map[string]int, error) {
	combinedData := make(map[string]int)
	for _, row := range records[1:] {
		if len(row) != len(ExpectedColumns) {
			return map[string]int{}, fmt.Errorf("wrong amount of data provided :: %v", row)
		}
		if row[0] == "" || row[1] == "" || row[2] == "" || row[3] == "" {
			return map[string]int{}, fmt.Errorf("empty data provided :: %v", row)
		}
		priceInt, err := strconv.Atoi(row[2])
		if err != nil {
			return map[string]int{}, fmt.Errorf("unable to convert :: got %v", row[2])
		}
		qtyInt, err := strconv.Atoi(row[3])
		if err != nil {
			return map[string]int{}, fmt.Errorf("unable to convert :: got %v", row[2])
		}
		combinedData[row[1]] = combinedData[row[1]] + priceInt*qtyInt
	}
	return combinedData, nil
}

func (r *RunAnalysis) checkValidQuick(records [][]string, expectedColumns []string) error {
	// Check empty
	if len(records) < 2 {
		return fmt.Errorf("empty dataset passed in")
	}
	// Check same number of columns
	if len(records[0]) != len(expectedColumns) {
		return fmt.Errorf("unexpected number of columns")
	}
	// Ensure csv header at the top is the same
	for i, h := range records[0] {
		if h != expectedColumns[i] {
			return fmt.Errorf("expected column name: %v, got: %v", expectedColumns[i], h)
		}
	}

	return nil
}

func (ra *RunAnalysis) downloadFile(ctx context.Context, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	o := ra.StorageSvc.Bucket(ra.BucketName).Object(fileName)
	rc, err := o.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to create storage reader :: %v", err)
	}
	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := rc.Close(); err != nil {
		return fmt.Errorf("reader.Close: %v", err)
	}
	return nil
}
