PROJECT_NAME=ecr-cleaner
VERSION=v0.5
GO_FILES := $(shell find ./ -name '*.go' | grep -v vendor)

all: bin                                                     ## Default rule

bin:                                                         ## Build binary
	GOARCH=amd64 GOOS=linux go build -o bin/linux/$(PROJECT_NAME)
	
clean:                                                       ## Clean code
	rm -rf bin

vet-check:                                                   ## Verify vet compliance
ifeq ($(shell go tool vet -all -shadow=true $(GO_FILES) 2>&1 | wc -l), 0)
	@printf "ok\tall files passed go vet\n"
else
	@printf "error\tsome files did not pass go vet\n"
	@go tool vet -all -shadow=true $(GO_FILES) 2>&1
endif

fmt-check:                                                   ## Verify fmt compliance
ifneq ($(shell gofmt -l $(GO_FILES) | wc -l), 0)
	@printf "error\tsome files did not pass go fmt: $(shell gofmt -l $(GO_FILES))\n"; exit 1
else
	@printf "ok\tall files passed go fmt\n"
endif

test:                                                        ## Test go code and coverage
	@go test -covermode=count -coverprofile=$(COVER_PROFILE) $(BINARY_PKG)

full-test: test fmt-check vet-check                          ## Pass test / fmt / vet

html-cover: test                                             ## Display coverage in HTML
	@go tool cover -html=$(COVER_PROFILE)

release: clean bin                                           ## Create release
	hub release create -a "bin/linux/$(PROJECT_NAME)" -m "$(VERSION)" "$(VERSION)"

package: clean bin                                           ## Create archive
	zip -j main.zip bin/linux/ecr-cleaner python/index.py

help:                                                        ## Show this help
	@printf "Rules:\n"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: bin clean test vet-check fmt-check test full-test html-cover release package help
