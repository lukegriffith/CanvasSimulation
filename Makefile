export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

.ONESHELL:

default: test build 


build:
	go build -i server ./cmd/server/main.go

test:
	-go test -fullpath ./...

