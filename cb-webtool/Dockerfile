FROM golang:1.11.2-alpine

WORKDIR /go/src/github.com/cloud-barista/cb-webtool 
COPY . .

RUN apk update && apk add git
RUN apk add --no-cache bash
RUN go get -u -v github.com/go-session/echo-session
RUN go get -u github.com/labstack/echo/... && go get github.com/cespare/reflex
RUN go get -u github.com/davecgh/go-spew/spew


EXPOSE 1234

CMD reflex -r '\.(html|go)' -s go run main.go
