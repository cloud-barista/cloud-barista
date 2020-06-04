#####################################################
##### STEP 1. Build executable binary
#####################################################

FROM golang:alpine AS builder 

# install git for fetching the dependencies
RUN apk update && apk add --no-cache git

WORKDIR /app

ADD go.mod .
ADD main.go .

# build the binary
RUN CGO_ENABLED=0 go build -o ./hmac_server -v ./...
RUN chmod +x ./hmac_server

##############################################################
## Stage 2 - Applicaiton Setup
##############################################################

FROM scratch

# set working directory and files
WORKDIR /app

COPY --from=builder /app/hmac_server .

ENTRYPOINT [ "/app/hmac_server" ]

EXPOSE 8010