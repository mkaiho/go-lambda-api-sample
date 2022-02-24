BIN_DIR:=bin
BIN_ZIP_DIR:=$(BIN_DIR)/zip
ROOT_PACKAGE:=github.com/mkaiho/go-lambda-api-sample
COMMAND_PACKAGES:=$(shell go list ./cmd/...)
BINARIES := $(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BIN_DIR)/%)
ZIP_BINARIES := $(BINARIES:$(BIN_DIR)/%=$(BIN_ZIP_DIR)/%)

.PHONY: build
build: clean $(BINARIES)

$(BINARIES): $(GO_FILES)
	@go build -o $@ $(@:$(BIN_DIR)/%=$(ROOT_PACKAGE)/cmd/%)

.PHONY: zip
zip: clean build $(BIN_ZIP_DIR) $(ZIP_BINARIES)
$(BIN_ZIP_DIR):
	@test -d $(BIN_ZIP_DIR) || mkdir $(BIN_ZIP_DIR)
$(ZIP_BINARIES): $(BINARIES)
	@zip -j $@.zip $(@:$(BIN_ZIP_DIR)/%=$(BIN_DIR)/%)

.PHONY: dev-deps
dev-deps:
	go get gotest.tools/gotestsum@v1.7.0
	go get github.com/vektra/mockery/v2/.../
	go mod tidy

.PHONY: deps
deps:
	go mod download

.PHONY: gen-mock
gen-mock:
	make dev-deps
	mockery --all --case underscore --recursive --keeptree

.PHONY: test
test:
	gotestsum ./entity/... ./usecase/... ./adapter/... ./infrastructure/...

.PHONY: test-report
test-report:
	rm -r ./test-results
	mkdir -p ./test-results
	gotestsum --junitfile ./test-results/unit-tests.xml ./entity/... ./usecase/... ./adapter/... ./infrastructure/...

.PHONY: clean
clean:
	@rm -rf ./bin
