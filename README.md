# Blind-msg
[![GoDoc](https://img.shields.io/badge/go-doc-green)](https://godoc.org/github.com/jinil-ha/blind-msg)
[![build](https://api.travis-ci.com/jinil-ha/blind-msg.svg?branch=master)](https://travis-ci.com/jinil-ha/blind-msg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jinil-ha/blind-msg)](https://goreportcard.com/report/github.com/jinil-ha/blind-msg)
[![license](https://img.shields.io/badge/license-coffee-blue)](https://github.com/jinil-ha/blind-msg/blob/master/LICENSE.md)

Blind LINE message service using QR Code.
Check https://blind-msg.jinstory.net/

## Getting Started

### Prerequistes
* nginx with https
* mysql
* LINE login channel
* LINE message channel
* slack bot

### Install
```sh
git clone https://github.com/jinil-ha/blind-msg
cd blind-msg
dep ensure
```

### Configuration
```sh
# config web server (ex. nginx)
## location for resource
location ~ ^/(img|css|js)/ {
  root  <repo-directory>/resource/;
}
## location for qr code
location /qr {
  root /home/user/download/qr;
}
## location blind-msg
location / {
  proxy_pass	http://127.0.0.1:10080/;
  proxy_set_header  X-REAL-IP $remote_addr;
}

# edit server.yaml
cp server.yaml-sample server.yaml
vi server.yaml
```

## Build & Run
```sh
# build
make

# run
make run

# or start daemon
make start
```
