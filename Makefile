docker: install
	docker build -t fosdem-2025/inventory-service -f ./build/inventory/Dockerfile .
	docker build -t fosdem-2025/pricing-service -f ./build/pricing/Dockerfile .
	docker build -t fosdem-2025/product-service -f ./build/product/Dockerfile .

install:
	GOBIN=$(shell pwd)/bin/ go install ./... 
