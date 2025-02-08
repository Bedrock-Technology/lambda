GO     := go
ENTRY  := *.go
DIST   := .
OUTPUT := ${DIST}/lambda

GO_VERSION   := $(shell go version | awk '{print $$3}')
GIT_DESCRIBE := $(shell git describe --always --tags --dirty)
GIT_HASH     := $(shell git rev-parse HEAD)
CURRENT_TIME := $(shell date +'%Y-%m-%d %H:%M:%S %z')
GO_OS        := $(shell go env GOOS)
GO_ARCH      := $(shell go env GOARCH)

GLOBAL_LD_FLAGS := -X 'main.Version=${GIT_DESCRIBE}' \
	-X 'main.GoVersion=${GO_VERSION}' \
	-X 'main.GitHash=${GIT_HASH}' \
	-X 'main.BuildTime=${CURRENT_TIME}'

.PHONY: build
build:
	@for _ in _ ; do \
		EXTRA_FLAGS="-X 'main.OSArch=${GO_OS}/${GO_ARCH}'" ; \
		${GO} build -ldflags "${GLOBAL_LD_FLAGS} $${EXTRA_FLAGS}" -o ${OUTPUT} ${ENTRY} ; \
	done

.PHONY: build-static
build-static:
	@for _ in _ ; do \
		EXTRA_FLAGS="-X 'main.OSArch=${GO_OS}/${GO_ARCH}'" ; \
		CGO_ENABLED=0 ${GO} build -ldflags "-extldflags -static ${GLOBAL_LD_FLAGS} $${EXTRA_FLAGS}" -o ${OUTPUT} ${ENTRY} ; \
	done
