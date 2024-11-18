LOCAL_BIN:=$(CURDIR)/bin


install-golangci-lint:
	GOBIN=${LOCAL_BIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	${LOCAL_BIN}/golangci-lint run ./... --config ./.golangci.pipeline.yaml

install-deps:
	GOBIN=${LOCAL_BIN} go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.2
