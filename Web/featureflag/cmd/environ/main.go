package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	tpl = `<!DOCTYPE html>
<html>
	<body>
		<h1>{{ .MainTitle }}</h1>
	</body>
</html>`
)

func main() {
	log.Println("begin server")
	exampleTitle := os.Getenv("TITLE")

	http.Handle("/", &index{Title: exampleTitle})
	s := &http.Server{
		Addr: "0.0.0.0:8080",
	}
	log.Fatal(s.ListenAndServe())
}

type index struct {
	Title string
}

func (i *index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start index handler")
	defer log.Println("end index handler")

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse"))
		return
	}

	data := struct {
		MainTitle string
	}{
		MainTitle: i.Title,
	}

	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse"))
		return
	}
}
