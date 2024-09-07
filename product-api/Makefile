swagger:
	@swagger generate spec -o swagger.yaml --scan-models

build:
	go build -o bin/microservices

run: build
	./bin/microservices
