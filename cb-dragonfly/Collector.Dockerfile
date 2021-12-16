###################################################
# Cloud-Barista CB-Dragonfly Module Dockerfile    #
###################################################

# Go 빌드 이미지 버전 및 알파인 OS 버전 정보
ARG BASE_IMAGE_BUILDER=golang
ARG GO_VERSION=1.15
ARG ALPINE_VERSION=3

###################################################
# 1. Build CB-Dragonfly binary file
###################################################

FROM ${BASE_IMAGE_BUILDER}:${GO_VERSION}-alpine AS go-builder

ENV CGO_ENABLED=1
ENV GO111MODULE="on"
ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GOPATH="/go"

#ARG GO_FLAGS="-mod=vendor"
ARG LD_FLAGS="-s -w"
ARG OUTPUT="bin/cb-dragonfly"

WORKDIR ${GOPATH}/src/github.com/cloud-barista/cb-dragonfly
COPY . ./

RUN apk add --update gcc

RUN apk add --no-cache \
    bash \
    build-base \
    gcc \
    make \
    musl-dev \
    tzdata \
    librdkafka-dev \
    pkgconf

RUN go build -tags musl ${GO_FLAGS} -ldflags "${LD_FLAGS}" -o ${OUTPUT} -i ./pkg/modules/procedure/push/collector/k8s_collector \
    && chmod +x ${OUTPUT}

###################################################
# 2. Set up CB-Dragonfly runtime environment
###################################################

FROM alpine:${ALPINE_VERSION} AS runtime-alpine

ENV TZ="Asia/Seoul"

RUN apk add --no-cache \
    bash \
    tzdata \
    librdkafka-dev \
    pkgconf \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime \
    && \
    echo "${TZ}" > /etc/timezone

###################################################
# 3. Execute CB-Dragonfly Module
###################################################

FROM runtime-alpine as cb-dragonfly
LABEL maintainer="innogrid <dev.cloudbarista@innogrid.com>"

ENV GOPATH="/go"
ENV CBSTORE_ROOT=${GOPATH}/src/github.com/cloud-barista/cb-dragonfly
ENV CBLOG_ROOT=${GOPATH}/src/github.com/cloud-barista/cb-dragonfly
ENV CBMON_ROOT=${GOPATH}/src/github.com/cloud-barista/cb-dragonfly

COPY --from=go-builder ${GOPATH}/src/github.com/cloud-barista/cb-dragonfly/file ${GOPATH}/src/github.com/cloud-barista/cb-dragonfly/file

WORKDIR /opt/cb-dragonfly
COPY --from=go-builder ${GOPATH}/src/github.com/cloud-barista/cb-dragonfly/bin/cb-dragonfly /opt/cb-dragonfly/bin/cb-dragonfly
RUN chmod +x /opt/cb-dragonfly/bin/cb-dragonfly \
    && ln -s /opt/cb-dragonfly/bin/cb-dragonfly /usr/bin

#EXPOSE 8094/udp
#EXPOSE 9090
#EXPOSE 9999

ENTRYPOINT ["cb-dragonfly"]
