GIT_UPDATE_INDEX    := $(shell git update-index --refresh)
GIT_VERSION         := $(shell git describe --tags --dirty 2>/dev/null || echo dev)
GO                  ?= $(shell which go)

ifneq ($(GO),)
GOENV_GOOS               := $(shell go env GOOS)
GOENV_GOARCH             := $(shell go env GOARCH)
GOENV_GOARM              := $(shell go env GOARM)
GOOS                     ?= $(GOENV_GOOS)
GOARCH                   ?= $(GOENV_GOARCH)
GOARM                    ?= $(GOENV_GOARM)
GO_BUILD_SRC             := $(shell find . -name \*.go -type f) go.mod go.sum
GO_BUILD_SRC             += $(shell find templates -type f)
GO_BUILD_SRC             += $(shell find static -type f)
GO_BUILD_EXTLDFLAGS      := -static
GO_BUILD_TAGS            := static
GO_BUILD_TARGET_DEPS     :=
GO_BUILD_FLAGS           :=
GO_BUILD_LDFLAGS_OPTIMS  :=

ifeq ($(GOOS)/$(GOARCH),$(GOENV_GOOS)/$(GOENV_GOARCH))
GO_BUILD_TARGET          := dist/abstraction.fr
GO_BUILD_VERSION_TARGET  := dist/abstraction.fr-$(GIT_VERSION)
else
GO_BUILD_TARGET          := dist/abstraction.fr-$(GOOS)-$(GOARCH)
GO_BUILD_VERSION_TARGET  := dist/abstraction.fr-$(GIT_VERSION)-$(GOOS)-$(GOARCH)
endif # $(GOOS)/$(GOARCH)

ifeq ($(shell uname),Linux)
GO_BUILD_EXTLDFLAGS      += -lbsd
endif

ifneq ($(DEBUG),)
GO_BUILD_FLAGS           += -race -gcflags="all=-N -l"
else
GO_BUILD_LDFLAGS_OPTIMS  += -s -w
endif # $(DEBUG)

GO_BUILD_EXTLDFLAGS      := $(strip $(GO_BUILD_EXTLDFLAGS))
GO_BUILD_TAGS            := $(strip $(GO_BUILD_TAGS))
GO_BUILD_TARGET_DEPS     := $(strip $(GO_BUILD_TARGET_DEPS))
GO_BUILD_FLAGS           := $(strip $(GO_BUILD_FLAGS))
GO_BUILD_LDFLAGS_OPTIMS  := $(strip $(GO_BUILD_LDFLAGS_OPTIMS))
GO_BUILD_LDFLAGS         := -ldflags '$(GO_BUILD_LDFLAGS_OPTIMS) -X main.version=$(GIT_VERSION) -extldflags "$(GO_BUILD_EXTLDFLAGS)"'
endif # $(GO)

GO_BUILD_FLAGS_TARGET                := .go-build-flags
GO_CROSSBUILD_PLATFORMS              ?= linux/amd64 windows/amd64 darwin/amd64 darwin/arm64
GO_CROSSBUILD_LINUX_PLATFORMS        := $(filter linux/%,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_FREEBSD_PLATFORMS      := $(filter freebsd/%,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_OPENBSD_PLATFORMS      := $(filter openbsd/%,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_WINDOWS_PLATFORMS      := $(filter windows/%,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_DARWIN_PLATFORMS       := $(filter darwin/%,$(GO_CROSSBUILD_PLATFORMS))
GO_CROSSBUILD_LINUX_TARGET_PATTERN   := dist/abstraction.fr-$(GIT_VERSION)-linux-%
GO_CROSSBUILD_FREEBSD_TARGET_PATTERN := dist/abstraction.fr-$(GIT_VERSION)-freebsd-%
GO_CROSSBUILD_OPENBSD_TARGET_PATTERN := dist/abstraction.fr-$(GIT_VERSION)-openbsd-%
GO_CROSSBUILD_WINDOWS_TARGET_PATTERN := dist/abstraction.fr-$(GIT_VERSION)-windows-%.exe
GO_CROSSBUILD_DARWIN_TARGET_PATTERN  := dist/abstraction.fr-$(GIT_VERSION)-darwin-%
GO_CROSSBUILD_TARGETS                := $(patsubst linux/%,$(GO_CROSSBUILD_LINUX_TARGET_PATTERN),$(GO_CROSSBUILD_LINUX_PLATFORMS))
GO_CROSSBUILD_TARGETS                += $(patsubst freebsd/%,$(GO_CROSSBUILD_FREEBSD_TARGET_PATTERN),$(GO_CROSSBUILD_FREEBSD_PLATFORMS))
GO_CROSSBUILD_TARGETS                += $(patsubst openbsd/%,$(GO_CROSSBUILD_OPENBSD_TARGET_PATTERN),$(GO_CROSSBUILD_OPENBSD_PLATFORMS))
GO_CROSSBUILD_TARGETS                += $(patsubst windows/%,$(GO_CROSSBUILD_WINDOWS_TARGET_PATTERN),$(GO_CROSSBUILD_WINDOWS_PLATFORMS))
GO_CROSSBUILD_TARGETS                += $(patsubst darwin/%,$(GO_CROSSBUILD_DARWIN_TARGET_PATTERN),$(GO_CROSSBUILD_DARWIN_PLATFORMS))

DOCKER_BUILD_IMAGE      ?= ghcr.io/sylr/abstraction.fr
DOCKER_BUILD_VERSION    ?= $(GIT_VERSION)
DOCKER_BUILD_GO_VERSION ?= 1.16rc1
DOCKER_BUILD_LABELS      = --label org.opencontainers.image.title=prometheus-azure-exporter
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.description="Azure metrics exporter for prometheus"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.url="https://github.com/sylr/abstraction.fr"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.source="https://github.com/sylr/abstraction.fr"
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.revision=$(GIT_REVISION)
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.version=$(GIT_VERSION)
DOCKER_BUILD_LABELS     += --label org.opencontainers.image.created=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
DOCKER_BUILD_BUILD_ARGS ?= --build-arg=GO_VERSION=$(DOCKER_BUILD_GO_VERSION)
DOCKER_BUILDX_PLATFORMS ?= "linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6"
DOCKER_BUILDX_CACHE     ?= /tmp/.buildx-cache
DOCKER_BUILD_TARGET     := .docker-build

# ------------------------------------------------------------------------------

.PHONY: all build

all: build

clean:
	@git clean -ndx
	@/bin/echo -n "Would you like to proceed (yes/no) ? "
	@read proceed && test "$$proceed" == "yes" && git clean -fdx
	@cd ./lib/librdkafka/ && git reset --hard

# -- tests ---------------------------------------------------------------------

.PHONY: test test-go

test: test-go

test-go:
	@go test ./...

# -- build ---------------------------------------------------------------------

.PHONY: build run build-go .FORCE

$(GO_BUILD_FLAGS_TARGET) : .FORCE
	@(echo "GO_VERSION=$(shell $(GO) version)"; \
	  echo "GO_GOOS=$(GOOS)"; \
	  echo "GO_GOARCH=$(GOARCH)"; \
	  echo "GO_GOARM=$(GOARM)"; \
	  echo "GO_BUILD_TAGS=$(GO_BUILD_TAGS)"; \
	  echo "GO_BUILD_FLAGS=$(GO_BUILD_FLAGS)"; \
	  echo 'GO_BUILD_LDFLAGS=$(subst ','\'',$(GO_BUILD_LDFLAGS))') \
	    | cmp -s - $@ \
	        || (echo "GO_VERSION=$(shell $(GO) version)"; \
	            echo "GO_GOOS=$(GOOS)"; \
	            echo "GO_GOARCH=$(GOARCH)"; \
	            echo "GO_GOARM=$(GOARM)"; \
	            echo "GO_BUILD_TAGS=$(GO_BUILD_TAGS)"; \
	            echo "GO_BUILD_FLAGS=$(GO_BUILD_FLAGS)"; \
	            echo 'GO_BUILD_LDFLAGS=$(subst ','\'',$(GO_BUILD_LDFLAGS))') > $@

run:
	$(GO) run .

build: build-go

scp: $(GO_BUILD_VERSION_TARGET)
	scp $(GO_BUILD_VERSION_TARGET) root@abstraction.fr:/usr/local/bin
	ssh root@abstraction.fr "systemctl stop abstraction.fr.service && \
		cd /usr/local/bin && \
		unlink abstraction.fr && \
		ln -s $(shell basename $(GO_BUILD_VERSION_TARGET)) abstraction.fr && \
		systemctl start abstraction.fr.service"

build-go: $(GO_BUILD_VERSION_TARGET) $(GO_BUILD_TARGET)

$(GO_BUILD_TARGET): $(GO_BUILD_VERSION_TARGET)
	@(test -e $@ && unlink $@) || true
	@ln $< $@

$(GO_BUILD_VERSION_TARGET): $(GO_BUILD_SRC) $(GO_GENERATE_TARGET) $(GO_BUILD_FLAGS_TARGET) | $(GO_BUILD_TARGET_DEPS)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -tags $(GO_BUILD_TAGS) $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

crossbuild: $(GO_BUILD_VERSION_TARGET) $(GO_CROSSBUILD_TARGETS)

$(GO_CROSSBUILD_LINUX_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(shell echo $* | cut -d '/' -f1) GOARM=$(shell echo $* | cut -d '/' -f2 | sed "s/^v//")) $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_FREEBSD_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	CGO_ENABLED=0 GOOS=freebsd GOARCH=$(shell echo $* | cut -d '/' -f1) GOARM=$(shell echo $* | cut -d '/' -f2 | sed "s/^v//") $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_OPENBSD_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	CGO_ENABLED=0 GOOS=openbsd GOARCH=$(shell echo $* | cut -d '/' -f1) GOARM=$(shell echo $* | cut -d '/' -f2 | sed "s/^v//") $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_WINDOWS_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	CGO_ENABLED=0 GOOS=windows GOARCH=$(shell echo $* | cut -d '/' -f1) GOARM=$(shell echo $* | cut -d '/' -f2 | sed "s/^v//") $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

$(GO_CROSSBUILD_DARWIN_TARGET_PATTERN): $(GO_BUILD_SRC) $(GO_BUILD_FLAGS_TARGET)
	CGO_ENABLED=0 GOOS=darwin GOARCH=$(shell echo $* | cut -d '/' -f1) GOARM=$(shell echo $* | cut -d '/' -f2 | sed "s/^v//") $(GO) build -tags $(GO_BUILD_TAGS),crossbuild $(GO_BUILD_FLAGS) $(GO_BUILD_LDFLAGS) -o $@

# -- tools ---------------------------------------------------------------------

.PHONY: git-hooks

git-hooks:
	@{ test -e contrib -a -e .git && cp contrib/git/hooks/pre-commit .git/hooks/ } || true

# -- docker --------------------------------------------------------------------

.PHONY: docker-build

docker-build:
	@docker buildx build . -f Dockerfilex \
		-t $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION) \
		--cache-to=type=local,dest=$(DOCKER_BUILDX_CACHE) \
		--platform=$(DOCKER_BUILDX_PLATFORMS) \
		$(DOCKER_BUILD_BUILD_ARGS) \
		$(DOCKER_BUILD_LABELS)

docker-push:
	@docker buildx build . -f Dockerfilex \
		--push \
		-t $(DOCKER_BUILD_IMAGE):$(DOCKER_BUILD_VERSION) \
		--cache-to=type=local,dest=$(DOCKER_BUILDX_CACHE) \
		--platform=$(DOCKER_BUILDX_PLATFORMS) \
		$(DOCKER_BUILD_BUILD_ARGS) \
		$(DOCKER_BUILD_LABELS)
