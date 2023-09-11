# gochatter
Simple chatting service written in Go

## Requirements
Create a `.ssh` under root directory (`gochatter/`). Create two private key and
cert. Requires `openssl` for macos.
```
cd gochatter
mkdir .ssh

openssl genpkey -algorithm RSA -out gochatter.key
openssl req -new -x509 -sha256 -key gochatter.key -out gochatter.crt -days 3650
```

## Installation
```
cd gochatter
go install ./cmd/gochatter
```

## Uses
```
gochatter server // runs gochatter server hosted locally
gochatter client // connects to locally hosted server
```
