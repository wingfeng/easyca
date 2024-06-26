FROM golang:1.16.3 as builder
WORKDIR ./
COPY ./ /app

ARG APP_VERSION
RUN export GO111MODULE="on"
RUN export GOPROXY=https://goproxy.io,direct

ENV CA_PORT=80
WORKDIR /app
RUN make APP_VERSION=${APP_VERSION}
#RUN go mod vendor
# RUN go build -mod=vendor -o easyca -a -ldflags "-extldflags '-static' -s -X 'main.AppName=${APP_NAME}' \
# 				 -X 'conf.AppVersion=${APP_VERSION}' \ 
# 				-X 'conf.BuildVersion=${BUILD_VERSION}' \
# 				-X 'conf.BuildTime=${BUILD_TIME}' \
# 				-X 'conf.GitRevision=${GIT_REVISION}' \
# 				-X 'conf.GitBranch=${GIT_BRANCH}' \
# 				-X 'conf.GoVersion=$GO_VERSION'" 



FROM alpine:latest
WORKDIR /
VOLUME ["/db","/certs"]
COPY --from=builder /app/easyca .
COPY /dist /dist
COPY /conf /conf
Run rm /conf/*.go

RUN /easyca -ver

#cluster port
EXPOSE 9000
CMD ["/easyca"]