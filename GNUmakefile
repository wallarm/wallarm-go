all: test

testacc:
	TF_ACC=1 go test ./... -test.v -timeout 120m

test:
	go test ./... -test.v -timeout=30s -parallel=4 -race

.PHONY: test testacc
