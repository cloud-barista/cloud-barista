#-------------------------------------------
# STEP 1 : build executable binary
#-------------------------------------------
FROM golang:1.16-alpine as builder

# gcc
RUN apk add --no-cache build-base

ADD . /usr/src/app

WORKDIR /usr/src/app

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -tags cb-ladybug -o cb-ladybug -v src/main.go

#-------------------------------------------
# STEP 2 : build a image
#-------------------------------------------
FROM scratch as prod

COPY --from=builder /usr/src/app/conf/ /app/conf/
COPY --from=builder /usr/src/app/cb-ladybug /app/
COPY --from=builder /usr/src/app/src/scripts/ /app/src/scripts/

ENV CBLOG_ROOT "/app"
ENV CBSTORE_ROOT "/app"
ENV APP_ROOT "/app"

ENV SPIDER_URL http://cb-spider:1024/spider
ENV TUMBLEBUG_URL http://cb-tumblebug:1323/tumblebug
ENV BASE_PATH /ladybug

ENV API_USERNAME default
ENV API_PASSWORD default

ENTRYPOINT [ "/app/cb-ladybug" ]

EXPOSE 8080
