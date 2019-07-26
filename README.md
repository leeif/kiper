# kiper
A small configuration tool which is a wrapper of [kingpin](https://github.com/alecthomas/kingpin) + [viper](https://github.com/spf13/viper.git). I like the way viper handles the config files however I am a fan of kinping because it can validate the command line flags easily. This library is just a combination of the features of this two projects.

# Feature

* Merge command line flags and config file settings.
* Validate command line flags and config file bothly.

# Usage

```
$ go get github.com/leeif/kiper
```

## Example

[kiper_example.go](https://github.com/leeif/kiper/blob/master/_example/kiper_example.go)

```
$ dep ensure
$ cd _example
$ go run kiper_example.go --help
```

## Struct Tags

kiper_config

|  Field  |  Description  |
| ---- | ---- |
|  name  |  config name  |

kiper_value

|  Field  |  Description  |
| ---- | ---- |
|  name  |  value name  |
|  default  |  default value  |
|  help  |  help message  |
