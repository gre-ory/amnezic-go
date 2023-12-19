VERSION=$(shell git describe --tags --always --dirty)
ifeq ($(version),)
	TAG=$(VERSION)
else
	TAG=$(version)
endif
BIN_DIR=bin
SERVER=${BIN_DIR}/amnezic-go

ifneq ($(verbose),)
	TEST_ARGS += -v
endif

ifeq ($(short),true)
	TEST_ARGS += -short
endif

PACKAGE=github.com/gre-ory/amnezic-go
# PACKAGE_CMD_DAEMON=${PACKAGE}/cmd/${SERVICE}

LDFLAGS = -X 'main.version=$(TAG)'

# run 'make Q="" <rule>' to enable verbosity
Q := @

.PHONY:	all build test install run

build:
	@echo " ----- build -----"
	$(Q)CGO_ENABLED=1 GOOS=linux go build $(GO_BUILD_FLAGS) -ldflags "${LDFLAGS}" -o ${SERVER} ${PACKAGE}
run: build db-up
	@echo " ----- run -----"
	@./scripts/run
test:
	$(Q) go test -race ./...
react:
	@./scripts/get-react.sh
push:
	@./scripts/push-docker.sh
db-status:
	@goose -dir "./db" "sqlite3" "./db/amnezic.db" status
db-up:
	@echo " ----- db-up -----"
	@goose -dir "./db" "sqlite3" "./db/amnezic.db" up
db-down:
	@goose -dir "./db" "sqlite3" "./db/amnezic.db" down
