# CB-MCKS
> Multi-Cloud Kubernetes Service Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/cloud-barista/cb-mcks)](https://goreportcard.com/report/github.com/cloud-barista/cb-mcks)
[![Build](https://img.shields.io/github/workflow/status/cloud-barista/cb-mcks/Build%20amd64%20container%20image)](https://github.com/cloud-barista/cb-mcks/actions?query=workflow%3A%22Build+amd64+container+image%22)
[![Top Language](https://img.shields.io/github/languages/top/cloud-barista/cb-mcks)](https://github.com/cloud-barista/cb-mcks/search?l=go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cloud-barista/cb-mcks?label=go.mod)](https://github.com/cloud-barista/cb-mcks/blob/master/go.mod)
[![Repo Size](https://img.shields.io/github/repo-size/cloud-barista/cb-mcks)](#)
[![GoDoc](https://godoc.org/github.com/cloud-barista/cb-mcks?status.svg)](https://pkg.go.dev/github.com/cloud-barista/cb-mcks@master)
[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-mcks?color=blue)](https://github.com/cloud-barista/cb-mcks/releases/latest)
[![License](https://img.shields.io/github/license/cloud-barista/cb-mcks?color=blue)](https://github.com/cloud-barista/cb-mcks/blob/master/LICENSE)

```
[NOTE]
CB-MCKS is currently under development. (The latest version is v0.5.0 (Affogato))
So, we do not recommend using the current release in production.
Please note that the functionalities of CB-MCKS are not stable and secure yet.
If you have any difficulties in using CB-MCKS, please let us know.
(Open an issue or Join the cloud-barista Slack)
```

## Getting started

### Preparation

* Golang 1.16.+ ([Download and install](https://golang.org/doc/install))

### Dependencies

* CB-Tumblebug [v0.5.0](https://github.com/cloud-barista/cb-tumblebug/releases/tag/v0.5.0)
* CB-Spider [v0.5.0](https://github.com/cloud-barista/cb-spider/releases/tag/v0.5.0)


### Clone

```
$ git clone https://github.com/cloud-barista/cb-mcks.git
$ cd cb-mcks
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
$ go build -o cb-mcks src/main.go
```

```
$ export CBLOG_ROOT="$(pwd)"
$ export CBSTORE_ROOT="$(pwd)"
$ nohup ./cb-mcks & > /dev/null
```

### Test

```
$ curl -s  http://localhost:1470/mcks/healthy -o /dev/null -w "code:%{http_code}"

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
* Open http://localhost:1470/swagger/index.html in your web browser 

## Documents

* [Design](./docs/design)
* REST API [latest](https://cloud-barista.github.io/cb-mcks-api-web/?url=https://raw.githubusercontent.com/cloud-barista/cb-mcks/master/src/docs/swagger.yaml)


## Contribution
Learn how to start contribution on the [Contributing Guideline](https://github.com/cloud-barista/docs/tree/master/contributing) and [Style Guideline](https://github.com/cloud-barista/cb-mcks/blob/master/STYLE_GUIDE.md)

