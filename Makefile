CGO_ENABLED?=0
allPkgs = $(shell go list ./...)

.PHONY: all
all: test static-analysis

.PHONY: test
test:
	go test -cover -mod=vendor ./...

.PHONY: static-analysis
static-analysis: lint vet errcheck

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

.PHONY: errcheck
errcheck:
	go run github.com/kisielk/errcheck -exclude errcheck-exclude.txt $(allPkgs)

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor

.PHONY: container
container: 
	docker build -t cbutera90/netdisco-exporter .

.PHONY: publish-container
publish-container: 
	docker push cbutera90/netdisco-exporter

.PHONY: run
run:
	docker run -p 8080:8080 -e NETDISCO_HOST -e NETDISCO_USERNAME -e NETDISCO_PASSWORD -it cbutera90/netdisco-exporter