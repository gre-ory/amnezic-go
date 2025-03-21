BIN_NAME=amnezic-go
VERSION=$(shell git describe --tags --always --dirty)
ifeq ($(version),)
	TAG=$(VERSION)
else
	TAG=$(version)
endif
BIN_DIR=bin
BIN=${BIN_DIR}/${BIN_NAME}

ifneq ($(verbose),)
	TEST_ARGS += -v
endif

ifeq ($(short),true)
	TEST_ARGS += -short
endif

PACKAGE=github.com/gre-ory/amnezic-go
PACKAGE_CMD=${PACKAGE}/cmd

LDFLAGS = -X 'main.version=$(TAG)'

# run 'make Q="" <rule>' to enable verbosity
Q := @

.PHONY:	all build test install run

build:
	@print-header "go build"
	$(Q)CGO_ENABLED=1 GOOS=linux go build $(GO_BUILD_FLAGS) -ldflags "${LDFLAGS}" -o ${BIN} ${PACKAGE_CMD}

pre-run: build db-up

run: pre-run
	@print-header "go-run"
	@./scripts/run.sh

loc: pre-run
	@print-header "run loc"
	@./scripts/run.sh loc

stg: pre-run
	@print-header "go-run stg"
	@./scripts/run.sh stg

prd: pre-run
	@print-header "go-run prd"
	@./scripts/run.sh prd

test:
	@print-header "go-test"
	$(Q) go test -race ./...

db-status:
	@print-header "db-status"
	@goose -dir "./db" "sqlite3" "./db/amnezic.db" status

db-up:
	@print-header "db up"
	@goose -dir "./db" "sqlite3" "./db/amnezic.db" up

db-down:
	@print-header "db-down"
	@goose -dir "./db" "sqlite3" "./db/amnezic.db" down
