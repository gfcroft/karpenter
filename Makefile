<<<<<<< HEAD
help: ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

presubmit: verify test licenses vulncheck ## Run all steps required for code to be checked in

test: ## Run tests
	go test ./... \
		-race \
		-timeout 10m \
		--ginkgo.focus="${FOCUS}" \
		--ginkgo.timeout=10m \
		--ginkgo.v \
		-cover -coverprofile=coverage.out -outputdir=. -coverpkg=./...

deflake: ## Run randomized, racing tests until the test fails to catch flakes
=======
export K8S_VERSION ?= 1.27.x
CLUSTER_NAME ?= $(shell kubectl config view --minify -o jsonpath='{.clusters[].name}' | rev | cut -d"/" -f1 | rev | cut -d"." -f1)

## Inject the app version into operator.Version
LDFLAGS ?= -ldflags=-X=sigs.k8s.io/karpenter/pkg/operator.Version=$(shell git describe --tags --always)

GOFLAGS ?= $(LDFLAGS)
WITH_GOFLAGS = GOFLAGS="$(GOFLAGS)"

## Extra helm options
CLUSTER_ENDPOINT ?= $(shell kubectl config view --minify -o jsonpath='{.clusters[].cluster.server}')
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)
KARPENTER_IAM_ROLE_ARN ?= arn:aws:iam::${AWS_ACCOUNT_ID}:role/${CLUSTER_NAME}-karpenter
HELM_OPTS ?= --set serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn=${KARPENTER_IAM_ROLE_ARN} \
      		--set settings.clusterName=${CLUSTER_NAME} \
			--set settings.interruptionQueue=${CLUSTER_NAME} \
			--set controller.resources.requests.cpu=1 \
			--set controller.resources.requests.memory=1Gi \
			--set controller.resources.limits.cpu=1 \
			--set controller.resources.limits.memory=1Gi \
			--set settings.featureGates.spotToSpotConsolidation=true \
			--create-namespace

# CR for local builds of Karpenter
KARPENTER_NAMESPACE ?= kube-system
KARPENTER_VERSION ?= $(shell git tag --sort=committerdate | tail -1)
KO_DOCKER_REPO ?= ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/dev
GETTING_STARTED_SCRIPT_DIR = website/content/en/preview/getting-started/getting-started-with-karpenter/scripts

# Common Directories
MOD_DIRS = $(shell find . -path "./website" -prune -o -name go.mod -type f -print | xargs dirname)
KARPENTER_CORE_DIR = $(shell go list -m -f '{{ .Dir }}' sigs.k8s.io/karpenter)

# TEST_SUITE enables you to select a specific test suite directory to run "make e2etests" or "make test" against
TEST_SUITE ?= "..."
TEST_TIMEOUT ?= "3h"

help: ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

presubmit: verify test ## Run all steps in the developer loop

ci-test: test coverage ## Runs tests and submits coverage

ci-non-test: verify licenses vulncheck ## Runs checks other than tests

run: ## Run Karpenter controller binary against your local cluster
	SYSTEM_NAMESPACE=${KARPENTER_NAMESPACE} \
		KUBERNETES_MIN_VERSION="1.19.0-0" \
		LEADER_ELECT=false \
		DISABLE_WEBHOOK=true \
		CLUSTER_NAME=${CLUSTER_NAME} \
		INTERRUPTION_QUEUE=${CLUSTER_NAME} \
		FEATURE_GATES="Drift=true" \
		go run ./cmd/controller/main.go

test: ## Run tests
	go test -v ./pkg/$(shell echo $(TEST_SUITE) | tr A-Z a-z)/... \
		-cover -coverprofile=coverage.out -outputdir=. -coverpkg=./... \
		--ginkgo.focus="${FOCUS}" \
		--ginkgo.randomize-all \
		--ginkgo.vv

e2etests: ## Run the e2e suite against your local cluster
	cd test && CLUSTER_ENDPOINT=${CLUSTER_ENDPOINT} \
		CLUSTER_NAME=${CLUSTER_NAME} \
		INTERRUPTION_QUEUE=${CLUSTER_NAME} \
		go test \
		-p 1 \
		-count 1 \
		-timeout ${TEST_TIMEOUT} \
		-v \
		./suites/$(shell echo $(TEST_SUITE) | tr A-Z a-z)/... \
		--ginkgo.focus="${FOCUS}" \
		--ginkgo.timeout=${TEST_TIMEOUT} \
		--ginkgo.grace-period=3m \
		--ginkgo.vv

e2etests-deflake: ## Run the e2e suite against your local cluster
	cd test && CLUSTER_NAME=${CLUSTER_NAME} ginkgo \
		--focus="${FOCUS}" \
		--timeout=${TEST_TIMEOUT} \
		--grace-period=3m \
		--until-it-fails \
		--vv \
		./suites/$(shell echo $(TEST_SUITE) | tr A-Z a-z) \

benchmark:
	go test -tags=test_performance -run=NoTests -bench=. ./...

deflake: ## Run randomized, racing, code-covered tests to deflake failures
	for i in $(shell seq 1 5); do make test || exit 1; done

deflake-until-it-fails: ## Run randomized, racing tests until the test fails to catch flakes
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	ginkgo \
		--race \
		--focus="${FOCUS}" \
		--timeout=10m \
		--randomize-all \
		--until-it-fails \
		-v \
		./pkg/...

<<<<<<< HEAD
=======
coverage:
	go tool cover -html coverage.out -o coverage.html

verify: tidy download ## Verify code. Includes dependencies, linting, formatting, etc
	go generate ./...
	hack/boilerplate.sh
	cp  $(KARPENTER_CORE_DIR)/pkg/apis/crds/* pkg/apis/crds
	hack/validation/requirements.sh
	hack/validation/labels.sh
	hack/github/dependabot.sh
	$(foreach dir,$(MOD_DIRS),cd $(dir) && golangci-lint run $(newline))
	@git diff --quiet ||\
		{ echo "New file modification detected in the Git working tree. Please check in before commit."; git --no-pager diff --name-only | uniq | awk '{print "  - " $$0}'; \
		if [ "${CI}" = true ]; then\
			exit 1;\
		fi;}
	@echo "Validating codegen/docgen build scripts..."
	@find hack/code hack/docs -name "*.go" -type f -print0 | xargs -0 -I {} go build -o /dev/null {}
	actionlint -oneline

>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
vulncheck: ## Verify code vulnerabilities
	@govulncheck ./pkg/...

licenses: download ## Verifies dependency licenses
	! go-licenses csv ./... | grep -v -e 'MIT' -e 'Apache-2.0' -e 'BSD-3-Clause' -e 'BSD-2-Clause' -e 'ISC' -e 'MPL-2.0'

<<<<<<< HEAD
verify: ## Verify code. Includes codegen, dependencies, linting, formatting, etc
	go mod tidy
	go generate ./...
	hack/validation/kubelet.sh
	hack/validation/taint.sh
	hack/validation/requirements.sh
	hack/validation/labels.sh
	@# Use perl instead of sed due to https://stackoverflow.com/questions/4247068/sed-command-with-i-option-failing-on-mac-but-works-on-linux
	@# We need to do this "sed replace" until controller-tools fixes this parameterized types issue: https://github.com/kubernetes-sigs/controller-tools/issues/756
	@perl -i -pe 's/sets.Set/sets.Set[string]/g' pkg/scheduling/zz_generated.deepcopy.go
	hack/boilerplate.sh
	go vet ./...
	golangci-lint run
	@git diff --quiet ||\
		{ echo "New file modification detected in the Git working tree. Please check in before commit."; git --no-pager diff --name-only | uniq | awk '{print "  - " $$0}'; \
		if [ "${CI}" = true ]; then\
			exit 1;\
		fi;}
	actionlint -oneline
=======
setup: ## Sets up the IAM roles needed prior to deploying the karpenter-controller. This command only needs to be run once
	CLUSTER_NAME=${CLUSTER_NAME} ./$(GETTING_STARTED_SCRIPT_DIR)/add-roles.sh $(KARPENTER_VERSION)

image: ## Build the Karpenter controller images using ko build
	$(eval CONTROLLER_IMG=$(shell $(WITH_GOFLAGS) KO_DOCKER_REPO="$(KO_DOCKER_REPO)" ko build --bare github.com/aws/karpenter-provider-aws/cmd/controller))
	$(eval IMG_REPOSITORY=$(shell echo $(CONTROLLER_IMG) | cut -d "@" -f 1 | cut -d ":" -f 1))
	$(eval IMG_TAG=$(shell echo $(CONTROLLER_IMG) | cut -d "@" -f 1 | cut -d ":" -f 2 -s))
	$(eval IMG_DIGEST=$(shell echo $(CONTROLLER_IMG) | cut -d "@" -f 2))

apply: image ## Deploy the controller from the current state of your git repository into your ~/.kube/config cluster
	helm upgrade --install karpenter charts/karpenter --namespace ${KARPENTER_NAMESPACE} \
        $(HELM_OPTS) \
        --set logLevel=debug \
        --set controller.image.repository=$(IMG_REPOSITORY) \
        --set controller.image.tag=$(IMG_TAG) \
        --set controller.image.digest=$(IMG_DIGEST)

install:  ## Deploy the latest released version into your ~/.kube/config cluster
	@echo Upgrading to ${KARPENTER_VERSION}
	helm upgrade --install karpenter oci://public.ecr.aws/karpenter/karpenter --version ${KARPENTER_VERSION} --namespace ${KARPENTER_NAMESPACE} \
		$(HELM_OPTS)

delete: ## Delete the controller from your ~/.kube/config cluster
	helm uninstall karpenter --namespace ${KARPENTER_NAMESPACE}

docgen: ## Generate docs
	KARPENTER_CORE_DIR=$(KARPENTER_CORE_DIR) $(WITH_GOFLAGS) ./hack/docgen.sh

codegen: ## Auto generate files based on AWS APIs response
	$(WITH_GOFLAGS) ./hack/codegen.sh

stable-release-pr: ## Generate PR for stable release
	$(WITH_GOFLAGS) ./hack/release/stable-pr.sh

snapshot: ## Builds and publishes snapshot release
	$(WITH_GOFLAGS) ./hack/release/snapshot.sh

release: ## Builds and publishes stable release
	$(WITH_GOFLAGS) ./hack/release/release.sh

release-crd: ## Packages and publishes a karpenter-crd helm chart
	$(WITH_GOFLAGS) ./hack/release/release-crd.sh

prepare-website: ## prepare the website for release
	./hack/release/prepare-website.sh

toolchain: ## Install developer toolchain
	./hack/toolchain.sh

issues: ## Run GitHub issue analysis scripts
	pip install -r ./hack/github/requirements.txt
	@echo "Set GH_TOKEN env variable to avoid being rate limited by Github"
	./hack/github/feature_request_reactions.py > "karpenter-feature-requests-$(shell date +"%Y-%m-%d").csv"
	./hack/github/label_issue_count.py > "karpenter-labels-$(shell date +"%Y-%m-%d").csv"

website: ## Serve the docs website locally
	cd website && npm install && hugo mod tidy && hugo server

tidy: ## Recursively "go mod tidy" on all directories where go.mod exists
	$(foreach dir,$(MOD_DIRS),cd $(dir) && go mod tidy $(newline))
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13

download: ## Recursively "go mod download" on all directories where go.mod exists
	$(foreach dir,$(MOD_DIRS),cd $(dir) && go mod download $(newline))

<<<<<<< HEAD
toolchain: ## Install developer toolchain
	./hack/toolchain.sh
=======
update-karpenter: ## Update kubernetes-sigs/karpenter to latest
	go get -u sigs.k8s.io/karpenter@HEAD
	go mod tidy
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539

<<<<<<< HEAD
.PHONY: help presubmit dev test verify toolchain
=======
.PHONY: help dev ci release test e2etests verify tidy download docgen codegen apply delete toolchain licenses vulncheck issues website nightly snapshot

define newline


endef
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
