.PHONY: init dep doc godoc mock lint lint-dupl test bench build clean all serve cov di nilaway fmt docker-build get-ssh-key

init:
	# pip install pre-commit
	# pre-commit install
	# go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8
	# for wire
	go install github.com/google/wire/cmd/wire@latest
	# for make doc
	go install github.com/swaggo/swag/cmd/swag@v1.16.2
	# for make mock
	go install github.com/golang/mock/mockgen@v1.6.0
	# for gofumpt
	go install mvdan.cc/gofumpt@latest
	# for golines
	go install github.com/segmentio/golines@latest

dep:
	go mod tidy

doc:
	swag init -o ./cmd/api/docs/ --parseDependency --parseInternal -g cmd/api/main.go

godoc:
	godoc -http=127.0.0.1:6060 -goroot="."

mock:
	go generate ./...

di:
	wire ./internal/di

lint:
	golangci-lint run

lint-dupl:
	golangci-lint run --no-config --disable-all --enable=dupl

nilaway:
	nilaway ./...

fmt:
	golines ./ -m 120 -w --base-formatter gofmt --no-reformat-tags
	gofumpt -l -w .

test:
	go test -gcflags=all=-l $(shell go list ./... | grep -v mock | grep -v docs) -covermode=count -coverprofile .coverage.cov

cov: test
	go tool cover -html=.coverage.cov

bench:
	go test -run=nonthingplease -benchmem -bench=. $(shell go list ./...)

build:
	CGO_ENABLED=0 go build -o bin/go-clean-arch cmd/api/main.go

all: lint test build

serve: build
	./bin/go-clean-arch -config configs/config-dev.yaml

clean:
	rm -rf bin
	rm -rf .coverage.cov

docker-build:
	docker build .
