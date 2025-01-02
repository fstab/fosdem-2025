package main

import (
	"encoding/json"
	"github.com/fstab/fosdem-2025/internal/model/inventory"
	"log"
	"net/http"
)

const usage = `<html><head><title>inventory service</title></head>
<body>
<h1>Inventory Service</h1>
Example query: <a href="/inventory?search=telescope">http://localhost:8081/inventory?search=telescope</a>
</body></html>
`

func searchHandler(w http.ResponseWriter, req *http.Request) {
	searchString := req.URL.Query().Get("search")
	if searchString == "" {
		defaultHandler(w, req)
		return
	}
	searchResult := search(searchString)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(searchResult)
	if err != nil {
		log.Printf("failed to send response: %s", err)
	}
}

func search(searchString string) []inventory.Item {
	return []inventory.Item{
		{1, "Celestron", 3},
		{2, "Meade", 7},
	}
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(usage))
	if err != nil {
		log.Printf("failed to send response: %s", err)
	}
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/inventory", searchHandler)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("error listening on port 8081: %s", err)
	}
}
