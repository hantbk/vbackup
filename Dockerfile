# syntax = docker/dockerfile:experimental
FROM node:14.18.2 AS buildvue

WORKDIR /vbackup/
COPY . /vbackup/
RUN --mount=type=cache,target=/vbackup/web/dashboard/node_modules cd /vbackup/web/dashboard &&\
    npm install
RUN --mount=type=cache,target=/vbackup/web/dashboard/node_modules cd /vbackup/web/dashboard &&\
    npm run build:prod

FROM golang:1.22.5-alpine3.20 AS buildbin
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPATH=/root/gopath

WORKDIR /vbackup/
COPY --from=buildvue /vbackup /vbackup/
RUN --mount=type=cache,target=/root/gopath apk update &&\
    apk upgrade &&\
    apk add --no-cache git make libffi-dev openssl-dev libtool tzdata curl &&\
    cp /usr/share/zoneinfo/UTC /etc/localtime &&  \
    sh prepare.sh
RUN --mount=type=cache,target=/root/gopath make build_go

FROM alpine:latest
LABEL MAINTAINER="vbackup <captainnemot1k60@gmail.com>"
ENV LANG C.UTF-8
COPY --from=buildbin /vbackup/dist/vbackup_server_* /apps/vbackup_server
COPY --from=buildbin /etc/localtime /etc/localtime

EXPOSE 8012

ENTRYPOINT ["/apps/vbackup_server"]


# # syntax = docker/dockerfile:experimental
# FROM node:14.18.2 AS buildvue

# WORKDIR /vbackup/
# COPY . /vbackup//
# RUN --mount=type=cache,target=/vbackup/web/dashboard/node_modules cd /vbackup/web/dashboard &&\
#     npm config set registry https://registry.npmmirror.com && npm install
# RUN --mount=type=cache,target=/vbackup/web/dashboard/node_modules cd /vbackup/web/dashboard &&\
#     npm run build:prod

# FROM golang:1.22.5-alpine3.20 AS buildbin
# ENV GO111MODULE=on
# # ENV GOPROXY="https://goproxy.cn,direct"
# ENV CGO_ENABLED=0
# ENV GOPATH=/root/gopath

# WORKDIR /vbackup/
# COPY --from=buildvue /vbackup /vbackup/
# RUN --mount=type=cache,target=/root/gopath echo -e 'https://mirrors.ustc.edu.cn/alpine/v3.20/main/\nhttps://mirrors.ustc.edu.cn/alpine/v3.20/community/' > /etc/apk/repositories &&\
#     apk update &&\
#     apk upgrade &&\
#     apk add --no-cache git make libffi-dev openssl-dev libtool tzdata curl &&\
#     cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&  \
#     sh prepare.sh
# RUN --mount=type=cache,target=/root/gopath make build_go

# FROM alpine:latest
# LABEL MAINTAINER="vbackup <captainnemot1k60@gmail.com>"
# ENV LANG C.UTF-8
# COPY --from=buildbin /vbackup/dist/vbackup_server_* /apps/vbackup_server
# COPY --from=buildbin /etc/localtime /etc/localtime

# EXPOSE 8012

# ENTRYPOINT ["/apps/vbackup_server"]

