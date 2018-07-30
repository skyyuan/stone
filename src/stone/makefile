all: build package push cleanup

export PACKAGE=stone
export REPO_PATH=/go/src/$(PACKAGE)
export BUILD_PATH=$(REPO_PATH)/build
export DOCKER_USER=`id -u $(USER)`

build:
	docker run --rm -u $(DOCKER_USER) -v $(PWD):$(REPO_PATH) golang:1.9.2 go build -ldflags "-X main.Version=$(TAG)" -o $(BUILD_PATH)/eth $(PACKAGE)/cmd/eth/
	docker run --rm -u $(DOCKER_USER) -v $(PWD):$(REPO_PATH) golang:1.9.2 go build -ldflags "-X main.Version=$(TAG)" -o $(BUILD_PATH)/web $(PACKAGE)/cmd/web/
	docker run --rm -u $(DOCKER_USER) -v $(PWD):$(REPO_PATH) golang:1.9.2 go build -ldflags "-X main.Version=$(TAG)" -o $(BUILD_PATH)/ethwatcher $(PACKAGE)/cmd/ethwatcher/
	docker run --rm -u $(DOCKER_USER) -v $(PWD):$(REPO_PATH) golang:1.9.2 go build -ldflags "-X main.Version=$(TAG)" -o $(BUILD_PATH)/chaos $(PACKAGE)/cmd/chaos/

package:
	docker build -t registry.cn-hangzhou.aliyuncs.com/wallet/blockchain:$(TAG) .

push:
	docker push registry.cn-hangzhou.aliyuncs.com/wallet/blockchain:$(TAG)

cleanup:
	rm -rf build/*
