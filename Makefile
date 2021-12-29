CMD = agent aggregator graph hbs judge nodata transfer gateway api alarm
TARGET = open-falcon
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
GOFMT ?= gofmt "-s"
VERSION := $(shell cat VERSION)

all: $(CMD) $(TARGET)

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

vet:
	go vet $(PACKAGES)

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
	if [ $@ = "gateway" ]; then \
		GO111MODULE=on go build -ldflags "-X main.BinaryName=gateway -X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" \
			-o bin/gateway/falcon-gateway ./modules/transfer ; \
	else \
		GO111MODULE=on go build -ldflags "-X main.BinaryName=$@ -X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" \
			-o bin/$@/falcon-$@ ./modules/$@ ; \
	fi

.PHONY: $(TARGET)
$(TARGET): $(GOFILES)
	GO111MODULE=on go build -ldflags "-X main.BinaryName=Open-Falcon -X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" -o open-falcon

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
	@(cd ./out && ln -s ./agent/public ./public)
	@(cd ./out && mkdir -p ./agent/plugin && ln -s ./agent/plugin ./plugin)
	@cp -r ./modules/api/data ./out/api/
	@mkdir out/graph/data
	@bash ./config/confgen.sh
	@cp $(TARGET) ./out/$(TARGET)
	tar -C out -zcf open-falcon-v$(VERSION).tar.gz .
	@rm -rf out

pack4docker: checkbin
	@if [ -e out ] ; then rm -rf out; fi
	@mkdir out
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/bin;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/config;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/logs;)
	@$(foreach var,$(CMD),cp ./config/$(var).json ./out/$(var)/config/cfg.json;)
	@$(foreach var,$(CMD),cp ./bin/$(var)/falcon-$(var) ./out/$(var)/bin;)
	@if expr "$(CMD)" : "agent" > /dev/null; then \
		(cp -r ./modules/agent/public ./out/agent/); \
		(cd ./out && ln -s ./agent/public ./public); \
		(cd ./out && mkdir -p ./agent/plugin && ln -s ./agent/plugin ./plugin); \
	fi
	@if expr "$(CMD)" : "api" > /dev/null; then \
		cp -r ./modules/api/data ./out/api/; \
	fi
	@if expr "$(CMD)" : "graph" > /dev/null; then \
		mkdir out/graph/data; \
	fi
	@cp ./docker/ctrl.sh ./out/ && chmod +x ./out/ctrl.sh
	@cp $(TARGET) ./out/$(TARGET)
	tar -C out -zcf open-falcon-v$(VERSION).tar.gz .
	@rm -rf out

.PHONY: test
test:
	@go test ./modules/api/test

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: clean all agent aggregator graph hbs judge nodata transfer gateway api alarm
