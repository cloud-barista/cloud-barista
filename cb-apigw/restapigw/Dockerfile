
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache bash git && apk --no-cache add ca-certificates

WORKDIR /app

ADD ./cmd ./cmd
ADD ./conf ./conf
ADD ./pkg ./pkg
ADD ./go.mod .
ADD ./main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -tags cb-restapigw -o cb-restapigw -v

#############################################################
## Stage 2 - Applicaiton Setup
##############################################################

FROM ubuntu:latest

# use bash
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cb-restapigw .
COPY --from=builder /app/conf ./conf

ENTRYPOINT [ "/app/cb-restapigw", "-c", "/app/conf/cb-restapigw.yaml", "-d", "run" ]

EXPOSE 8000

# Image Build
# docker build -t cb-restapigw .
# Run Container
# docker run -it --name cb-restapigw -p 8000:8000 --network deploy_default cb-restapigw
