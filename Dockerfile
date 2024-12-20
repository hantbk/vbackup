# syntax = docker/dockerfile:experimental
FROM node:14.18.2 AS buildvue

WORKDIR /vbackup/
COPY . /vbackup/
RUN --mount=type=cache,target=/vbackup/web/dashboard/node_modules cd /vbackup/web/dashboard &&\
    npm install
RUN --mount=type=cache,target=/vbackup/web/dashboard/node_modules cd /vbackup/web/dashboard &&\
    npm run build:prod

FROM golang:1.23-alpine3.20 AS buildbin
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPATH=/root/gopath

WORKDIR /vbackup/
COPY --from=buildvue /vbackup /vbackup/
RUN --mount=type=cache,target=/root/gopath apk update &&\
    apk update && \
    apk upgrade && \
    apk add --no-cache \
        git \
        make \
        libffi-dev \
        openssl-dev \
        libtool \
        tzdata \
        curl && \
    cp /usr/share/zoneinfo/UTC /etc/localtime && \
    sh prepare.sh
RUN --mount=type=cache,target=/root/gopath make build_go

FROM alpine:latest
ENV LANG C.UTF-8
COPY --from=buildbin /vbackup/dist/vbackup_server_* /apps/vbackup_server
COPY --from=buildbin /etc/localtime /etc/localtime

EXPOSE 8012

ENTRYPOINT ["/apps/vbackup_server"]