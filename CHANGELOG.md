# Changelog

## [0.5.1] (2024-01-29)

### Improve

* In interactive mode, the value of the zk node supports strings containing spaces.

## [0.5.0] (2024-01-26)

### New Features

* Added command `deleteall`, which will recursively delete the specified node and all its child nodes.

### Build/Testing/Packaging Improvement

* Support MacOSX arm64 architecture in release-all.

## [0.4.0] (2021-02-08)

### New Features

* Support configuration file:
  * zkcli will try to read ~/.config/zkcli.conf as default configuration(can be changed with `-config` argument) if the file exists.

### Changes

* No longer output logs from zk library. You can run command with `-v` argument to output logs.

### Internal changes

* Switch from unmaintained samuel/go-zookeeper/zk library to the new official
  upstream go-zookeeper/zk, version 1.0.2. Fixes connection and authentication
  bugs.

## [0.3.0] (2018-09-16)

### New Features

* Supports scrollbar when there are too many matched suggestions

### Improve

* Revert #7: Ignored with the suffix slash of zpath when completion
* Remove / from end of path when fetch data from zookeeper

### Internal changes

* Change to use go mod
* Upgrade go-prompt to v0.2.2 and go-zookeeper to v0.0.0-20180130194729-c4fab1ac1bec


## [0.2.0] (2018-05-30)

### New Features

* Add `-version` to show version info ([add69127e](https://github.com/maskshell/zkcli/commit/add69127e15a855f934629ef437286d416122fc8))

```
$ zkcli -version
Version:	0.2.0
Git commit:	9fd746b
Built: 2018-05-30T13:44:21+0000
```

### Internal changes

* Remove unnecessary qiniupkg ([1c33d63f590](https://github.com/maskshell/zkcli/commit/1c33d63f590598c166ef0fcb4eb6554ca8bdee1c))
* Close connection before exit ([4c5d6a4d](https://github.com/maskshell/zkcli/commit/4c5d6a4dc16d28deec01df6c873e69b27b985f61))
* Ignored with the suffix slash of zpath when completion ([#7](https://github.com/let-us-go/zkcli/pull/7))


## 0.1.0 (2017-12-23)

* Initial Release


[0.2.0]: https://github.com/maskshell/zkcli/compare/v0.1.0...v0.2.0
[0.3.0]: https://github.com/maskshell/zkcli/compare/v0.2.0...v0.3.0
[0.4.0]: https://github.com/maskshell/zkcli/compare/v0.3.0...v0.4.0
[0.4.1]: https://github.com/maskshell/zkcli/compare/v0.4.0...v0.4.1
