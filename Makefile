export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

.PHONY: bin
bin:  
	mkdir -p $(GOBIN)
	cd tools && go mod tidy
  
.PHONY: bin.golangci-lint  
bin.golangci-lint: bin
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
  
.PHONY: lint
lint: bin.golangci-lint
	$(GOBIN)/golangci-lint run --timeout 3m
