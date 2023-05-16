IMG ?= qwwebhook:v1
NAME ?= qwwebhook

##@ Help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: build
build: ## 构建二进制文件
	CGO_ENABLED=0 go build -v -o bin/qwwebhook cmd/main.go

.PHONY: build-image
build-image: build ## 构建镜像
	docker build -t ${IMG} .
	docker save ${IMG} -o ${IMG}.tar
	ansible k8s -m synchronize -a "src=./${IMG}.tar dest=/tmp/"	
	ansible k8s -m shell -a "docker load -i /tmp/${IMG}.tar"
	rm -rf ${IMG}.tar

##@ TLS
.PHONY: tls
tls: ## 创建 webhook tls
	chmod +x bin/create_tls.sh
	./bin/create_tls.sh

##@ Install

.PHONY: install
install: ## 部署 webhook 资源
	kubectl apply -f ./manifests/
