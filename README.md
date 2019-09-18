A URL shortener written in GoLang. Works with any SQL variant.   
Supports configuration of key length, url normalization, collisions handling etc 

## Installation
* Install Golang 1.10+
* Install MySql
* Setup `GOPATH`, `GOBIN` path variables
* Install [Dep](https://github.com/golang/dep) to manage dependencies
* `./build.sh`

## Running binary:
Build script will generate two binaries - one for linux amd64 arch and other for OSX.   
Run right binary with the configuration file [example](https://github.com/AnkurGel/duss/blob/master/configs/config.yaml) as:    
`DUSS_CONFIG_PATH=/path/to/config.yaml ./releases/duss_linux_amd64`

## Development
* `go get -u github.com/ankurgel/duss`
* `cd $GOPATH/src/github.com/ankurgel/duss`
* `dep ensure`
* Edit `configs/config.yml` for MySql credentials
* Create relevant database
* `go run cmd/duss/main.go`

## Roadmap
[ ] - Domain blacklisting   
[ ] - Handle redirects and loops   
[ ] - Handle expiry of redirect   
[ ] - Support other data stores - Bolt, Redis   
