# gochatter
Simple chatting service written in Go

## Requirements
For now you'll need to retrieve a client secret from the creator. Once you have
the secret string, run the commands below.
```
cd gochatter
mkdir .secrets
echo "secret" > .secrets/clientsecret
```

## Installation
```
cd gochatter
make
```

## Uses
```
gochatter server // runs gochatter server hosted locally
gochatter client // connects to locally hosted server
```
