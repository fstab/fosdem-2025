package main

import (
	"encoding/json"
	"fmt"
	"github.com/fstab/fosdem-2025/internal/model/inventory"
	"github.com/fstab/fosdem-2025/internal/model/pricing"
	"github.com/fstab/fosdem-2025/internal/model/product"
	"log"
	"net/http"
)

const inventory_service_url = "http://inventory-service:8081"
const pricing_service_url = "http://pricing-service:8082"

const usage = `<html><head><title>product service</title></head>
<body>
<h1>Product Service</h1>
Example query: <a href="/products?search=telescope">http://localhost:8080/products?search=telescope</a>
</body></html>
`

func searchHandler(w http.ResponseWriter, req *http.Request) {
	searchString := req.URL.Query().Get("search")
	if searchString == "" {
		defaultHandler(w, req)
		return
	}
	inventoryItems, err := searchInventory(searchString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to query the inventory serice: %s", err)
		return
	}
	result := make([]product.Item, 0, len(inventoryItems))
	for _, item := range inventoryItems {
		price, err := queryPrice(item.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to query the pricing serice: %s", err)
			return
		}
		result = append(result, product.Item{
			Id:    item.Id,
			Name:  item.Name,
			Price: price,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("failed to send response: %s", err)
	}
}

func searchInventory(searchString string) ([]inventory.Item, error) {
	inventoryItems := []inventory.Item{}
	url := fmt.Sprintf("%s/inventory?search=%s", inventory_service_url, searchString)
	if err := queryJsonData(url, &inventoryItems); err != nil {
		return nil, err
	}
	return inventoryItems, nil
}

func queryPrice(productId int) (float64, error) {
	price := pricing.Price{}
	url := fmt.Sprintf("%s/prices/%d", pricing_service_url, productId)
	fmt.Println(url)
	if err := queryJsonData(url, &price); err != nil {
		return 0, err
	}
	return price.Price, nil
}

// Run an HTTP GET request and deserialize the JSON response into v.
func queryJsonData(url string, v any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with HTTP status %s", resp.Status)
	}
	responseContentType := resp.Header.Get("Content-Type")
	if responseContentType != "application/json" {
		if responseContentType == "" {
			return fmt.Errorf("expected content-type: application/json but received a response without a content-type header")
		} else {
			return fmt.Errorf("expected content-type: application/json but received content-type: %v", responseContentType)
		}
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("received illegal JSON response: %s", err)
	}
	return nil
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(usage))
	if err != nil {
		log.Printf("failed to send response: %s", err)
	}
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/products", searchHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error listening on port 8080: %s", err)
	}
}
