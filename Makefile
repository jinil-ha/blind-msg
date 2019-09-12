EXENAME := blind-msg
BASE := $(shell pwd)

GOEXEC := go
GOBIN := $(GOPATH)/bin
TIMEOUT  = 15

.PHONY: build
build: dep
	@cd ${BASE} && ${GOEXEC} build -v -o $(GOBIN)/$(EXENAME)

.PHONY: dep
dep: Gopkg.lock
	@cd ${BASE} && dep ensure

.PHONY: start
start: build
	@${GOBIN}/$(EXENAME) start

.PHONY: stop
stop:
	@${GOBIN}/$(EXENAME) stop

.PHONY: run
run: build
	@${GOBIN}/$(EXENAME) run

.PHONY: test
test:
	@cd ${BASE} && ${GOEXEC} test ./... -v
