ROOT_PACKAGE:=github.com/mkaiho/go-lambda-api-sample
BIN_DIR:=bin
SRC_DIR:=$(shell go list ./cmd/... ./lambda/...)
BINARIES:=$(SRC_DIR:$(ROOT_PACKAGE)/%=$(BIN_DIR)/%)
ARCHIVE_DIR:=$(BIN_DIR)/zip
ARCHIVES:=$(SRC_DIR:$(ROOT_PACKAGE)/%=$(ARCHIVE_DIR)/%)

.PHONY: build
build: clean $(BINARIES)

$(BINARIES):
	go build  -o $@ $(@:$(BIN_DIR)/%=$(ROOT_PACKAGE)/%)

.PHONY: archive
archive: $(ARCHIVES)

$(ARCHIVES):$(BINARIES)
	@test -d $(ARCHIVE_DIR) || mkdir $(ARCHIVE_DIR)
	@test -d $(ARCHIVE_DIR)/lambda || mkdir $(ARCHIVE_DIR)/lambda
	@test -d $(ARCHIVE_DIR)/cmd || mkdir $(ARCHIVE_DIR)/cmd
	@zip -j $@.zip $(@:$(ARCHIVE_DIR)/%=$(BIN_DIR)/%)
	@zip -j $@.zip jwks.json

.PHONY: dev-deps
dev-deps:
	go install gotest.tools/gotestsum@v1.7.0
	go install github.com/vektra/mockery/v2@latest

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
	@rm -rf ./test-results
	@mkdir -p ./test-results
	gotestsum --junitfile ./test-results/unit-tests.xml ./entity/... ./usecase/... ./adapter/... ./infrastructure/...

.PHONY: clean
clean:
	@rm -rf ./bin
