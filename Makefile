GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS)
VERSION=$(shell cat VERSION)

GoVersion=$(shell go version)
BASEPATH=$(shell pwd)
BUILD_TIME=$(shell date +"%Y%m%d%H%M")
BUILDDIR=$(BASEPATH)/dist
DASHBOARDDIR=$(BASEPATH)/web/dashboard
MAIN=$(BASEPATH)/cmd/main.go
APPVERSION=$(VERSION)

APP_NAME=vbackup_server_$(APPVERSION)_$(GOOS)_$(GOARCH)

LDFLAGS=-ldflags "-s -w -X github.com/hantbk/vbackup.BuildTime=${BUILD_TIME} -X github.com/hantbk/vbackup.V=${VERSION}"


all: build_web_dashboard all_bin
	$(BASEPATH)/hashsum.sh

all_bin: clean build_linux_amd64 build_linux_arm64 build_osx_amd64 build_osx_arm64

clean:
	rm -rf $(BUILDDIR)

# Build the web dashboard
build_web_dashboard:
	cd $(DASHBOARDDIR) && npm config set registry https://registry.npmmirror.com && npm install && npm run build:prod

build_go:
	go mod download
	go mod tidy
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath $(LDFLAGS) -o $(BUILDDIR)/$(APP_NAME) $(MAIN)

# Build binary files and web dashboard
build_bin: build_web_dashboard clean build_go

build_linux_amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath $(LDFLAGS) -o $(BUILDDIR)/vbackup_server_$(APPVERSION)_linux_amd64 $(MAIN)

build_linux_arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -trimpath $(LDFLAGS) -o $(BUILDDIR)/vbackup_server_$(APPVERSION)_linux_arm64 $(MAIN)

build_osx_amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -trimpath $(LDFLAGS) -o $(BUILDDIR)/vbackup_server_$(APPVERSION)_darwin_amd64 $(MAIN)

build_osx_arm64:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -trimpath $(LDFLAGS) -o $(BUILDDIR)/vbackup_server_$(APPVERSION)_darwin_arm64 $(MAIN)

# Build Docker image
build_image:
	docker buildx build -t vbackup/vbackup:${VERSION} -t vbackup/vbackup:latest --platform=linux/arm64,linux/amd64 . --push
