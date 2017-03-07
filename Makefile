SHELL := /bin/bash
TARGET_SOURCE = $(shell find main.go g cmd common -name '*.go')
CMD = agent aggregator graph hbs judge nodata query sender task transfer gateway api alarm
TARGET = open-falcon

VERSION := $(shell cat VERSION)

all: trash $(CMD) $(TARGET)

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
	@cp -r ./modules/agent/public ./out/agent/bin
	@cp -r ./modules/api/data ./out/api/
	@cp -r ./modules/alarm/views ./out/alarm/bin
	@mkdir out/graph/data
	@bash ./config/confgen.sh
	@cp $(TARGET) ./out/$(TARGET)
	tar -C out -zcf open-falcon-v$(VERSION).tar.gz .
	#@rm -rf out

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf ./package_cache_tmp
	@rm -rf ./vendor
	@rm -rf open-falcon-v$(VERSION).tar.gz

trash:
	go get -u github.com/rancher/trash
	trash -k -cache package_cache_tmp

.PHONY: trash clean all agent aggregator graph hbs judge nodata query sender task transfer gateway api alarm
