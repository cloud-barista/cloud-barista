# CB-Ladybug :beetle:
> Multi-Cloud Application Management Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/cloud-barista/cb-ladybug)](https://goreportcard.com/report/github.com/cloud-barista/cb-ladybug)
[![Build](https://img.shields.io/github/workflow/status/cloud-barista/cb-ladybug/Build%20amd64%20container%20image)](https://github.com/cloud-barista/cb-ladybug/actions?query=workflow%3A%22Build+amd64+container+image%22)
[![Top Language](https://img.shields.io/github/languages/top/cloud-barista/cb-ladybug)](https://github.com/cloud-barista/cb-ladybug/search?l=go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cloud-barista/cb-ladybug?label=go.mod)](https://github.com/cloud-barista/cb-ladybug/blob/master/go.mod)
[![Repo Size](https://img.shields.io/github/repo-size/cloud-barista/cb-ladybug)](#)
[![GoDoc](https://godoc.org/github.com/cloud-barista/cb-ladybug?status.svg)](https://pkg.go.dev/github.com/cloud-barista/cb-ladybug@master)
[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-ladybug?color=blue)](https://github.com/cloud-barista/cb-ladybug/releases/latest)
[![License](https://img.shields.io/github/license/cloud-barista/cb-ladybug?color=blue)](https://github.com/cloud-barista/cb-ladybug/blob/master/LICENSE)

```
[NOTE]
CB-Ladybug is currently under development. (the latest release is v0.4.0)
So, we do not recommend using the current release in production.
Please note that the functionalities of CB-Ladybug are not stable and secure yet.
If you have any difficulties in using CB-Ladybug, please let us know.
(Open an issue or Join the cloud-barista Slack)
```

## Getting started

### Preparation

* Golang 1.16.+ ([Download and install](https://golang.org/doc/install))

### Dependencies

* CB-Tumblebug [v0.4.0](https://github.com/cloud-barista/cb-tumblebug/releases/tag/v0.4.0)
* CB-Spider [v0.4.0](https://github.com/cloud-barista/cb-spider/releases/tag/v0.4.0)


### Clone

```
$ git clone https://github.com/cloud-barista/cb-ladybug.git
$ cd cb-ladybug
$ go get -v all
```

### Run 

```
$ export CBLOG_ROOT="$(pwd)"
$ export CBSTORE_ROOT="$(pwd)"
$ go run src/main.go
```

### Build and Execute

```
$ go build -o cb-ladybug src/main.go
```

```
$ export CBLOG_ROOT="$(pwd)"
$ export CBSTORE_ROOT="$(pwd)"
$ nohup ./cb-ladybug & > /dev/null
```

### Test

```
$ curl -s  http://localhost:8080/ladybug/healthy -o /dev/null -w "code:%{http_code}"

code:200
```


### API documentation

* Execute or Run
* Generate an updated swagger.yaml
```
$ go get -u github.com/swaggo/swag/cmd/swag

# in src/ folder
$ swag init
```
* Open http://localhost:8080/swagger/index.html in your web browser 

## Documents

* [Design](./docs/design)
* REST API [latest](https://cloud-barista.github.io/cb-ladybug-api-web/?url=https://raw.githubusercontent.com/cloud-barista/cb-ladybug/master/src/docs/v0.4.0.yaml)


## Contribution
Learn how to start contribution on the [Contributing Guideline](https://github.com/cloud-barista/docs/tree/master/contributing) and [Style Guideline](https://github.com/cloud-barista/cb-ladybug/blob/master/STYLE_GUIDE.md)

