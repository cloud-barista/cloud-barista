#############################################################
## Stage 1 - Client Build (Yarn)
#############################################################

# FROM mhart/alpine-node:12 AS client_builder
FROM node:10-alpine AS client_builder

WORKDIR /app

COPY ./web .

# Admin Web Build with yarn
RUN cd /app \
    && yarn install --pure-lockfile \
    && yarn generate

# CMD [ "ls -a /app "]

#############################################################
## Stage 2 - Golang Build (client embedding 포함)
#############################################################

FROM golang:1.15-alpine AS builder

RUN apk update && \
    apk add --no-cache bash git && \
    apk --no-cache add ca-certificates

# download static source builder
RUN go get -u github.com/mjibson/esc 

WORKDIR /app

# Copy sources
COPY ./cmd ./cmd
COPY ./conf ./conf
COPY ./pkg ./pkg
COPY ./go.mod .
COPY ./main.go .

# Copy Admin Web generated files
COPY --from=client_builder /app/dist ./web/dist 

# Admin web embedding to go source
RUN esc -o ./pkg/static/static.go -ignore "DS_Store|.nojekyll|README.md" -pkg static ./web/dist

# build go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -tags cb-restapigw -o cb-restapigw -v

#############################################################
## Stage 3 - Applicaiton Setup
#############################################################

FROM ubuntu:latest

# use bash
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cb-restapigw .
COPY --from=builder /app/conf ./conf

# Set ROOT path for CB-LOG, CB-STORE
ENV CBLOG_ROOT=/app CBSTORE_ROOT=/app

ENTRYPOINT [ "/app/cb-restapigw", "-c", "/app/conf/cb-restapigw.yaml", "-d", "run" ]

EXPOSE 8000 8001

# Image Build
# docker build -t cb-restapigw .

# Run Container
# docker run -it --rm --name cb-restapigw -p 8000:8000 -p 8001:8001 --network deploy_default cb-restapigw

# Delete stopped container
# docker rm cb-restapigw
