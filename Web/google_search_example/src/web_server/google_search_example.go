package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
	"log"
	"time"

	"search"
)

// handleSearch handles URLs like "/search?q=golang"
// Can take on varying query params like json and prettyjson
func handleSearch(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	// Check query
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, `missing "q" URL paramer`, http.StatusBadRequest)
		return
	}

	start := time.Now()
	//results, err := search.Search(query)
	//results, err := search.SearchParallel(query)
	results, err := search.SearchTimeout("golang", 80*time.Millisecond)
	elapsed := time.Since(start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Result []search.Result
		Elapsed time.Duration
	}

	resp := response{results, elapsed}

	switch req.FormValue("output"){
		case "json": 
			err = json.NewEncoder(w).Encode(resp)
		case "prettyjson":
			var b []byte
			b, err = json.MarshalIndent(resp, "","  ")
			if err == nil {
				_, err = w.Write(b)
			}
		default:
			err = responseTemplate.Execute(w, resp)
	}

	if err != nil {
		log.Print(err)
		return
	}
}

var responseTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
  <ol>
  {{range .Results}}
    <li>{{.Title}} - <a href="{{.URL}}">{{.URL}}</a></li>
  {{end}}
  </ol>
  <p>{{len .Results}} results in {{.Elapsed}}</p>
</body>
</html>
`))



func main() {
	http.HandleFunc("/search", handleSearch)
	fmt.Println("serving on http://localhost:8080/search")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}