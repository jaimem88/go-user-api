SHELL := /bin/bash

# Service specific vars
export COMPONENT_NAME := go-user-api
export COMPONENT_ARGS := $(if $(COMPONENT_ARGS),$(COMPONENT_ARGS),)

export SERVICE_PORT := $(if $(SERVICE_ADDR),$(SERVICE_ADDR),8899)
export SERVICE_ADDR := $(if $(SERVICE_ADDR),$(SERVICE_ADDR),:$(SERVICE_PORT))

export GOLANG_IMAGE :=devlube/gobuilder:0.0.36-msp-go-1.13.4-alpine3.10
export GOLANG_LINTER_IMAGE :=golangci/golangci-lint:v1.21.0

export GOFLAGS:=-mod=vendor

.PHONY: proto
proto:
	./pkg/proto/gen-go.sh

.PHONY: run
run: build
	docker build -t ${COMPONENT_NAME}-image .
	docker run -u 1000:1000 -p ${SERVICE_PORT}:${SERVICE_PORT} ${COMPONENT_NAME}-image

.PHONY: build
build:
	docker run --rm \
		-v ${PWD}:/usr/src/${COMPONENT_NAME} -w /usr/src/${COMPONENT_NAME} \
		-e GOFLAGS=${GOFLAGS} \
		-e GOOS=linux -e GOARCH=amd64 \
		${GOLANG_IMAGE} go build -ldflags="-w -s" -o ./bin/$(COMPONENT_NAME) ./cmd/$(COMPONENT_NAME)/

.PHONY: lint
lint:
	docker run --rm \
		-v ${PWD}:/go/src/app -w /go/src/app \
		${GOLANG_LINTER_IMAGE} golangci-lint run ${GOLANG_LINTER_ARGS}

.PHONY: test-all
test-all: test test-e2e

.PHONY: test
test:
	docker run --rm \
		-v ${PWD}:/usr/src/${COMPONENT_NAME} -w /usr/src/${COMPONENT_NAME} \
		-e GOFLAGS=${GOFLAGS} \
		${GOLANG_IMAGE} go test ./... -timeout 30s -count=1 ${ARGS}

.PHONY: test-e2e
test-e2e: build
	docker run --rm \
	-v ${PWD}:/go/app \
	-e COMPONENT_NAME=${COMPONENT_NAME} \
	-e GOFLAGS=${GOFLAGS} \
	-e SERVICE_ADDR=${SERVICE_ADDR} \
	-w /go/app \
	${GOLANG_IMAGE} \
	bash -c "./e2e.sh"

.PHONY: build-cov test-cov test-e2e-cov cov-all
test-cov:
	docker run --rm \
		-v ${PWD}:/usr/src/${COMPONENT_NAME} -w /usr/src/${COMPONENT_NAME} \
		-e GOFLAGS=${GOFLAGS} \
		${GOLANG_IMAGE} go test ./... -timeout 30s -count=1 -coverpkg=./internal/... -covermode=atomic -coverprofile=unit-cov.out ${ARGS}
build-cov:
	docker run --rm \
		-v ${PWD}:/usr/src/${COMPONENT_NAME} -w /usr/src/${COMPONENT_NAME} \
		-e GOFLAGS=${GOFLAGS} \
		${GOLANG_IMAGE} go test -c -o ./bin/${COMPONENT_NAME}-with-cov -cover -covermode=atomic -coverpkg=./internal/...,./pkg/... -count=1 ./cmd/${COMPONENT_NAME}

test-e2e-cov: build-cov
	docker run --rm \
	-v ${PWD}:/go/app \
	-e COMPONENT_NAME=${COMPONENT_NAME} \
	-e GOFLAGS=${GOFLAGS} \
	-e SERVICE_ADDR=${SERVICE_ADDR} \
	-w /go/app \
	${GOLANG_IMAGE} \
	bash -c "./e2e-cov.sh"

cov-all: test-cov test-e2e-cov cov-report

cov-report:
	docker run --rm \
		-v ${PWD}:/usr/src/${COMPONENT_NAME} -w /usr/src/${COMPONENT_NAME} \
		${GOLANG_IMAGE} bash -c "ls -la | grep .out && go get github.com/wadey/gocovmerge && \
			gocovmerge unit-cov.out e2e-cov.out > all-cov.out && \
			echo Unit Test coverage: && \
			go tool cover -func=unit-cov.out | grep total && \
			echo E2E Test coverage: && \
			go tool cover -func=e2e-cov.out | grep total && \
			echo Overall coverage: && \
			go tool cover -func=all-cov.out | grep total"
