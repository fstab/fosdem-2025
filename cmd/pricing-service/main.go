package main

import (
	"encoding/json"
	"github.com/fstab/fosdem-2025/internal/model/pricing"
	"github.com/fstab/fosdem-2025/internal/util"
	"log"
	"math/rand"
	"net/http"
)

const usage = `<html><head><title>pricing service</title></head>
<body>
<h1>Pricing Service</h1>
Example query: <a href="/prices/2">http://pricing-service:8082/prices/2</a>
</body></html>
`

func pricingHandler(w http.ResponseWriter, req *http.Request) {
	util.Sleep()
	if rand.Float64() < 0.05 { // simulate 5% error rate
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}
	productId := req.PathValue("productId")
	if productId == "" {
		defaultHandler(w, req)
		return
	}
	price := pricing.Price{
		ProductID: productId,
		Price:     float64(rand.Intn(1000)) / 100.0,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(price)
	if err != nil {
		log.Printf("failed to send response: %s", err)
	}
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	util.Sleep()
	_, err := w.Write([]byte(usage))
	if err != nil {
		log.Printf("failed to send response: %s", err)
	}
}

func main() {
	http.HandleFunc("/prices/{productId}", pricingHandler)
	http.HandleFunc("/", defaultHandler)
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("error listening on port 8081: %s", err)
	}
}
