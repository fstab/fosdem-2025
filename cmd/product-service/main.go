package main

import (
	"encoding/json"
	"fmt"
	"github.com/fstab/fosdem-2025/internal/model/inventory"
	"github.com/fstab/fosdem-2025/internal/model/pricing"
	"github.com/fstab/fosdem-2025/internal/model/product"
	"github.com/fstab/fosdem-2025/internal/util"
	"log"
	"net/http"
	"strings"
)

const inventory_service_url = "http://inventory-service:8081"
const pricing_service_url = "http://pricing-service:8082"

const usage = `<html><head><title>product service</title></head>
<body>
<h1>Product Service</h1>
Example query: <a href="/products?search=telescope">http://localhost:8080/products?search=telescope</a>
</body></html>
`

// handler for the /products?search=... endpoint
func searchHandler(w http.ResponseWriter, req *http.Request) {
	util.Sleep()
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
	result, err := queryPricesInParallel(inventoryItems)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to query the pricing serice: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("failed to send response: %s", err)
	}
}

// query the pricing service for each inventory item
func queryPricesInParallel(inventoryItems []inventory.Item) ([]product.Item, error) {
	prices := make(chan *pricing.Price)
	errors := make(chan error)
	for _, item := range inventoryItems {
		go func() {
			price, err := queryPrice(item.Id)
			if err != nil {
				errors <- err
			} else {
				prices <- price
			}
		}()
	}
	result := make([]product.Item, 0, len(inventoryItems))
	for i := 0; i < len(inventoryItems); i++ {
		select {
		case err := <-errors:
			return nil, err
		case price := <-prices:
			item := findById(price.ProductID, inventoryItems)
			if item == nil {
				return nil, fmt.Errorf("pricing service returned a product with an unknown product id: %s", price.ProductID)
			}
			result = append(result, product.Item{
				Id:    item.Id,
				Name:  item.Name,
				Price: price.Price,
			})
		}
	}
	return result, nil
}

func findById(id string, items []inventory.Item) *inventory.Item {
	for _, item := range items {
		if strings.Compare(fmt.Sprintf("%d", item.Id), id) == 0 {
			return &item
		}
	}
	return nil
}

// query the pricing service for a single inventory item
func queryPrice(productId int) (*pricing.Price, error) {
	price := pricing.Price{}
	url := fmt.Sprintf("%s/prices/%d", pricing_service_url, productId)
	if err := queryJsonData(url, &price); err != nil {
		return nil, err
	} else {
		return &price, nil
	}
}

// GET request to the inventory service
func searchInventory(searchString string) ([]inventory.Item, error) {
	inventoryItems := []inventory.Item{}
	url := fmt.Sprintf("%s/inventory?search=%s", inventory_service_url, searchString)
	if err := queryJsonData(url, &inventoryItems); err != nil {
		return nil, err
	}
	return inventoryItems, nil
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
	util.Sleep()
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
