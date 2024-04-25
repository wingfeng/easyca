
BINARY="easyca"
OUTPUT_DIR=.
GOOS=$(shell go env GOOS)
#GOOS
APP_NAME="企业证书管理系统"

BUILD_TIME=$(shell date "+%FT%T%z")
GIT_REVISION=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git name-rev --name-only HEAD)
APP_VERSION=0.0.1
BUILD_VERSION=$(APP_VERSION).$(GIT_REVISION)
GO_VERSION=$(shell go version)
release:	
	cd frontend
	yarn build
	cd ..	
	go env -w CGO_ENABLED=1;go build -a -installsuffix cgo -v  \
	-ldflags "-extldflags '-static' -s -X 'main.AppName=${APP_NAME}' \
				-X 'main.AppVersion=${APP_VERSION}' \
				-X 'main.BuildVersion=${BUILD_VERSION}' \
				-X 'main.BuildTime=${BUILD_TIME}' \
				-X 'main.GitRevision=${GIT_REVISION}' \
				-X 'main.GitBranch=${GIT_BRANCH}' \
				-X 'main.GoVersion=${GO_VERSION}'" \
	-o ${OUTPUT_DIR}/${BINARY} .
	./easyca -ver
docker:
	docker build . -t easyca:dev --build-arg APP_NAME=${APP_NAME} --build-arg APP_VERSION=${APP_VERSION} --build-arg BUILD_TIME="${BUILD_TIME}" --build-arg	 GIT_REVISION=${GIT_REVISION} --build-arg GIT_BRANCH=${GIT_BRANCH}