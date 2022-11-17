# Repo/Branch/Release
MODULE := $(shell basename `pwd`)
COMMIT := $(shell git log --pretty=format:'%h' -n 1)
TAG := $(shell git for-each-ref --count=1 --format='%(refname:short)' 'refs/tags/v[0-9]*.[0-9]*.[0-9]*' --points-at master --merged)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Structure
SRC := src
CMD_PATH := cmd/${MODULE}/*.go
BUILD_PATH := build
DEPLOY_PATH := deploy

# env
include ${BUILD_PATH}/ci/docker-compose/.env
export

# Package
# the conditional variable assignment operator, it can be changed during build for different envs

# Get latest merged tag in master, to allow release. Else, get the branch name as version and skip tags in there
VERSION ?= ""
ifeq ($(VERSION),"")
	ifeq ($(TAG),)
		VERSION    			= $(BRANCH)
	else
		VERSION				= $(TAG)
	endif
endif


# DOCKER_TAG from version var, but DNS compliant
DOCKER_TAG :=$(shell echo $(VERSION) | $(UBUNTU_CMD) awk '{gsub("[^.0-9a-zA-Z]","-");print $$0}' )
DOCKER_IMAGE_REPO ?= ""
KUBECONFIG ?= ""

# Build
GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(shell which go)
TIMESTAMP := $(shell date '+%Y-%m-%d_%I:%M')
LD_FLAGS = "-s -w -X $(MODULE)/pkg/version.BuildVersion=$(VERSION) -X $(MODULE)/pkg/version.BuildHash=$(COMMIT) -X $(MODULE)/pkg/version.BuildTime=$(TIMESTAMP)"
GOBUILD_OPTS = -ldflags=${LD_FLAGS}
GO_VERSION :=$(shell $(UBUNTU_CMD) awk '/^go / {print $$2}' $(SRC)/go.mod )


# Bins
GOLANG_CI_SHA := sha256:94388e00f07c64262b138a7508f857473e30fdf0f59d04b546a305fc12cb5961
CHART_TESTING_SHA := sha256:ef453de0be68d5ded26f3b3ea0c5431b396c8c48f89e2a07be7b19c4c9a68b31
HELM_SHA := sha256:6b85088a38ef34bbbdf3b91ab4e18038f35220f0f1bb1a97f94b7fde50ce66ee
GO_WORKSPACE_SHA := sha256:6b494c932ee8c209631e27521ddbe364da56e7f1275998fbb182447d20103e46
UBUNTU_SHA := sha256:f0a63f53b736b9211a5313a7219f6cc012b7cf4194c7ce2248fac8162b56dceb

GOCI_CMD := docker run --rm\
		-u $(id -u):$(id -g)\
		-v $(PWD)/${SRC}/:/opt/${SRC}\
		-w /opt/${SRC}\
		golangci/golangci-lint@${GOLANG_CI_SHA}

CT_CONTAINER_CMD := docker run -it --network host\
		-u $(id -u):$(id -g)\
		-v $(PWD)/${BUILD_PATH}/:/opt/${BUILD_PATH}\
		-v $(PWD)/${DEPLOY_PATH}/:/opt/${DEPLOY_PATH}\
		-v $(PWD)/.git/:/opt/.git:ro\
		-w "/opt"\
		quay.io/helmpack/chart-testing@${CHART_TESTING_SHA}

GO_WORKSPACE_CMD := docker run -i --rm\
		-u $(id -u):$(id -g)\
		-v $(PWD)/$(SRC):/go/$(MODULE)/$(SRC):ro\
		-w "/go/$(MODULE)/$(SRC)"\
		-e CGO_ENABLED=0\
		golang@${GO_WORKSPACE_SHA}


UBUNTU_CMD := docker run -i --rm\
		-u $(id -u):$(id -g)\
		ubuntu@${UBUNTU_SHA}


ifeq ($(KUBECONFIG),"")
	HELM_CONTAINER_CMD:=docker run --rm\
			-u $(id -u):$(id -g)\
			-v $(PWD)/${DEPLOY_PATH}/:/opt/${DEPLOY_PATH}:ro\
			-v ~/.kube:/root/.kube:ro\
			-w "/opt/${DEPLOY_PATH}"\
			alpine/helm@${HELM_SHA}
else
	HELM_CONTAINER_CMD:=docker run --rm\
			-u $(id -u):$(id -g)\
			-v $(PWD)/${DEPLOY_PATH}/:/opt/${DEPLOY_PATH}:ro\
			-v $(PWD)/${KUBECONFIG}:/root/.kube:ro\
			-w "/opt/${DEPLOY_PATH}"\
			alpine/helm@${HELM_SHA}
endif

HELP_CMD:=awk '{\
					if ($$0 ~ /^.PHONY: [a-zA-Z\-\_0-9]+$$/) {\
						command = substr($$0, index($$0, ":") + 2);\
						if (info) {\
							printf "\t\033[36m%-20s\033[0m %s\n",\
								command, info;\
							info = "";\
						}\
					} else if ($$0 ~ /^[a-zA-Z\-\_0-9.]+:/) {\
						command = substr($$0, 0, index($$0, ":"));\
						if (info) {\
							printf "\t\033[36m%-20s\033[0m %s\n",\
								command, info;\
							info = "";\
						}\
					} else if ($$0 ~ /^\#\#/) {\
						if (info) {\
							info = info"\n\t\t\t     "substr($$0, 3);\
						} else {\
							info = substr($$0, 3);\
						}\
					} else {\
						if (info) {\
							print "\n"info;\
						}\
						info = "";\
					}\
				}'				

##Local tools


## makefile iterate it self line by line and prints help
## in order to add info to your command, append "## <text>" one before the target
.PHONY: help
help:
	@cat $(firstword $(MAKEFILE_LIST)) | $(UBUNTU_CMD) $(HELP_CMD)

##Local Development

## Prints the current tag,branch and version
.PHONY: environment
environment:
	@echo "Tag: "${TAG}
	@echo "Branch: "${BRANCH} 
	@echo "Version: "${VERSION}
	@echo "Image Tag: "${DOCKER_TAG}
	@echo "Go path: "${GOPATH}
	@echo "Go bin: "${GOBIN}
	@echo "Go Version: "${GO_VERSION}


## Runs go mod {tidy,vendor,verify}
mod-sync: 	
	@cd ${SRC} && ${GOBIN} mod tidy && go mod vendor && go mod verify && echo "at: `pwd`"

## Runs local DB Server
.PHONY: local-infra
local-infra: 
	@docker-compose -f ${BUILD_PATH}/ci/docker-compose/local-infra.yaml up -d


## Runs local DB and the service for local development.
## Exports env vars of docker-compose
## Call it with optional cli arguments
## e.g. make local-run CMD_ARGS=query
.PHONY: local-run
local-run: local-infra
	@cd ${SRC} && ${GOBIN} run -mod=vendor ${GOBUILD_OPTS} ${CMD_PATH} $(CMD_ARGS)

## Detects the default exposed port from container's image, and run the container with the exposed port
docker-run: 
	{ 	\
		port=$$(docker inspect --format='{{range $$key, $$value := .Config.ExposedPorts }}{{$$key}}{{end}}' ${MODULE}:${VERSION} | sed 's/\/.*//') ;\
		docker run --env-file $(DEPLOY_PATH)/docker/.env --rm -p 127.0.0.1:$${port}:$${port} $(MODULE):$(VERSION) ;\
	}

##Continuous Integration

## Runs linter
lint:
	@echo "Linting...\n"
	@$(GOCI_CMD) golangci-lint run

## Runs tests
test:
	@echo "Unit Testing...\n"
	@$(GO_WORKSPACE_CMD) go test ./...

## Runs hadolint on Dockerfile
.PHONY: docker-lint
docker-lint: 
	@echo "Dockerfile linting..."
	@docker run --rm -i hadolint/hadolint < ${BUILD_PATH}/package/Dockerfile

## Builds image. Call it with VERSION arg to parse Image tag. 
## e.g. `make docker-image VERSION=feat/packaging_dockerfile`
docker-image: 
	@echo "Bulding image with tag: ${DOCKER_TAG}..."
	@DOCKER_BUILDKIT=0 docker build -t ${MODULE}:${DOCKER_TAG}  --build-arg LD_FLAGS=$(LD_FLAGS) -f ${BUILD_PATH}/package/Dockerfile ${SRC}

## Runs ct linting
.PHONY: chart-testing
chart-testing:
	@echo "Chart testing..."
	$(CT_CONTAINER_CMD) ct lint --config ${BUILD_PATH}/ci/ct/chart-testing.yaml --all

## Runs lints code and dockerfile, tests code and chart, builds docker image
ci:	lint test docker-lint chart-testing docker-image
	@echo "Continues Integration finished..."

##Continuous Delivery

## Tags and Pushes image to a Registry.
## Please set your repository and authenticate before calling this target!
.PHONY: cd
cd: ## Tags and Pushes image to a Registry. Currently to DockerHub
	@echo "Warning!! You must authenticate with DockerHub, and have repo access"
	@echo "Tagging ${MODULE}:${DOCKER_TAG} ${DOCKER_IMAGE_REPO}:${DOCKER_TAG}"
	@docker tag ${MODULE}:${DOCKER_TAG} ${DOCKER_IMAGE_REPO}:${DOCKER_TAG}
	@echo "Pushing ${DOCKER_IMAGE_REPO}:${DOCKER_TAG}"
	@docker push ${DOCKER_IMAGE_REPO}:${DOCKER_TAG}


##Continuous Deployment

## Optional args VERSION and KUBECONFIG. 
## e.g make deploy-prod VERSION=v1.1.7 KUBECONFIG=kube
.PHONY: deploy
deploy:
	$(HELM_CONTAINER_CMD) upgrade --install --create-namespace --namespace $(MODULE) $(MODULE) $(DEPLOY_PATH)/helm-chart/$(MODULE) --set image.tag=${DOCKER_TAG}

##Pipelines

## Lints, tests, builds and distributes
ci-cd: ci cd
	@echo "CI/CD finished...."

## Lints, tests, builds, distributes and deploys
ci-cd-deploy: ci cd deploy
	@echo "CI/CD finished...."


##Release

## Builds linux bin for amd64 arch
build-linux: 
	cd ${SRC} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=readonly ${GOBUILD_OPTS} -o ../${BUILD_PATH}/${MODULE} ${CMD_PATH}

## Builds darwin bin for amd64 arch
build-darwin:
	cd ${SRC} && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -mod=readonly ${GOBUILD_OPTS} -o ../${BUILD_PATH}/${MODULE}-darwin ${CMD_PATH}

## Builds all bins
build-all: build-linux build-darwin
	sha256sum ${BUILD_PATH}/${MODULE} > ${BUILD_PATH}/${MODULE}.sha256
	sha256sum ${BUILD_PATH}/${MODULE}-darwin> ${BUILD_PATH}/${MODULE}-darwin.sha256
