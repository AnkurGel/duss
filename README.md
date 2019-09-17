## Installation:
* Install Golang 1.10+
* Install MySql
* Setup `GOPATH`
* Install [Dep](https://github.com/golang/dep) to manage dependencies
* `go get -u http://github.com/AnkurGel/duss`
* `cd $GOPATH/src/github.com/AnkurGel/duss`
* `dep ensure`
* Edit `configs/config.yml` for MySql credentials
* Create relevant database
* `go run cmd/duss/main.go`
