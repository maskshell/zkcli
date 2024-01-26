# zkcli

[![Go Report Card](https://goreportcard.com/badge/github.com/maskshell/zkcli)](https://goreportcard.com/report/github.com/maskshell/zkcli)

A interactive Zookeeper client.

![zkcli](./zkcli.gif)


## Install

### Mac (Homebrew)

For original edition from [@let-us-go](https://github.com/let-us-go/zkcli).

```
brew tap let-us-go/zkcli
brew install zkcli
```

### go install

```
go install github.com/maskshell/zkcli@latest
```

### Build

```
make release-all
```

Or [download a pre-built binary](https://github.com/maskshell/zkcli/releases) for Linux or macOS/OSX.


## Usage

```shell
$ zkcli ls /test
[abc]
```

```shell
$ zkcli
>>> 
>>> help
get <path>
ls <path>
create <path> [<data>]
set <path> [<data>]
delete <path>
deleteall <path>
connect <host:port>
addauth <scheme> <auth>
close
exit
>>>
```

```shell
$ # zkcli will try to read ~/.config/zkcli.conf as default configuration if the file exists
$ zkcli -h # show help
```


## Develop

```
$ GO111MODULE=on go get xxx

$ GO111MODULE=on go mod tidy
$ GO111MODULE=on go mod verify

$ make build
$ make release-all
```
