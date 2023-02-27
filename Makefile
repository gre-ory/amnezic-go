VERSION=$(shell git describe --tags --always --dirty)
ifeq ($(version),)
	TAG=$(VERSION)
else
	TAG=$(version)
endif
BIN_DIR=bin
SERVER=${BIN_DIR}/server

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
	$(Q)CGO_ENABLED=0 GOOS=linux go build $(GO_BUILD_FLAGS) -ldflags "${LDFLAGS}" -o ${SERVER} ${PACKAGE}
run: build
	ENVIRONMENT="dev" APPLICATION_NAME="amnezic" APPLICATION_VERSION="9.9.9" LOG_CONFIG="dev" LOG_LEVEL="info" FRONTEND_ADDRESS=":9090" BACKEND_ADDRESS=":9091" ${SERVER}
test:
	$(Q) go test -race ./...
