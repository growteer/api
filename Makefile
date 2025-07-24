all:
	@CGO_ENABLED=0 go build -ldflags '-s -w' -o growteer-api ./cmd/growteer-api/...

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@mkdir -p reports
	@go test -v -covermode=atomic -coverprofile=reports/coverage.out ./... | tee reports/test-result.out
	@go-junit-report < reports/test-result.out > reports/junit-report.xml
	@go tool cover -html=reports/coverage.out -o reports/coverage.html

.PHONY: gen
gen:
	@gqlgen generate
	@mockery
	@go generate ./...
