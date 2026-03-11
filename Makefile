.PHONY: generate test vet

generate:
	controller-gen object paths=./api/...

test:
	go test ./...

vet:
	go vet ./...
