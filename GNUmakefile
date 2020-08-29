all: test

test:
	go test ./... -v -timeout=30s -parallel=4 -race

vet:
	go vet $(go list ./...)

cover:
	go test ./... -coverprofile c.out
	go tool cover -func c.out
	go tool cover -html=c.out -o coverage.html

.PHONY: test vet cover
