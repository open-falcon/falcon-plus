SHELL := /bin/bash
TARGET_SOURCE = $(shell find main.go g cmd common -name '*.go')
CMD = agent aggregator graph hbs judge nodata transfer gateway api alarm
TARGET = open-falcon
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
GOFMT ?= gofmt "-s"

VERSION := $(shell cat VERSION)

all: install $(CMD) $(TARGET)

install:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	@if [ -f ~/.bash_profile ]; then source ~/.bash_profile; fi
	govendor sync

fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

$(CMD):
	go get ./modules/$@
	go build -o bin/$@/falcon-$@ ./modules/$@

$(TARGET): $(TARGET_SOURCE)
	go build -ldflags "-X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" -o open-falcon

checkbin: bin/ config/ open-falcon

pack: checkbin
	@if [ -e out ] ; then rm -rf out; fi
	@mkdir out
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/bin;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/config;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/logs;)
	@$(foreach var,$(CMD),cp ./config/$(var).json ./out/$(var)/config/cfg.json;)
	@$(foreach var,$(CMD),cp ./bin/$(var)/falcon-$(var) ./out/$(var)/bin;)
	@cp -r ./modules/agent/public ./out/agent/
	@(cd ./out && ln -s ./agent/public/ ./public)
	@cp -r ./modules/agent/plugins ./out/agent/
	@(cd ./out && ln -s ./agent/plugins/ ./plugins)
	@cp -r ./modules/api/data ./out/api/
	@mkdir out/graph/data
	@bash ./config/confgen.sh
	@cp $(TARGET) ./out/$(TARGET)
	tar -C out -zcf open-falcon-v$(VERSION).tar.gz .
	@rm -rf out

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: clean all agent aggregator graph hbs judge nodata transfer gateway api alarm
