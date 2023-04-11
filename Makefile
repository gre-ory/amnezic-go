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
	$(Q)CGO_ENABLED=0 GOOS=linux go build $(GO_BUILD_FLAGS) -ldflags "${LDFLAGS}" -o ${SERVER} ${PACKAGE}
run: build
	@./scripts/run.sh
test:
	$(Q) go test -race ./...
react:
	@./scripts/get-react.sh
push:
	@./scripts/push-docker.sh
