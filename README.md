# Kiper
Although there are lots of solutions for handling command line flags or config file, sometimes we need to handle the command line flags and config file at the same time. And kiper is aiming at providing a convenience at this scenario. 

Rewriting the command line flags and config file parser is reinventing the wheel, so the kiper is just a wrapper of [kingpin](https://github.com/alecthomas/kingpin) + [viper](https://github.com/spf13/viper.git). I like the way viper handles the config files however personally I am a fan of kinping because it's fluent-style and it can validate the command line flags easily. This library is a combination of this two projects to handle the flags and config files at the same time.

## Feature

* Support various type of data like int, string, bool, struct, array(slice)
* Merge flag and config file settings automatically.
* Custom validations for both flag and config file settings.

## install

```
$ go get github.com/leeif/kiper
```

## Tags

*kiper_config*: an entry point of the kiper parser. kiper_config is a collection of kiper_values.

|  Field  |  Description  |
| ---- | ---- |
|  name  |  config name  |

*kiper_value*: the actual data filed you need to parse from the command line flags and config file.

|  Field  |  Description  |
| ---- | ---- |
|  name  |  value name  |
|  default  |  default value  |
|  help  |  help message  |
|  required      | set the flag to required, default value will be ignored due to the kingpin feature |

## KiperValue interface

You can implement the KiperValue interface int your own data structs to achieve the data validation.

```
type KiperValue interface {
	Set(string) error
	String() string
}
```

```
type Port struct {
	p string
}

func (port *Port) Set(p string) error {
	if p == "" {
		return errors.New("port can't be empty")
	}
	port.p = p
	return nil
}

func (port *Port) String() string {
	return port.p
}

type Config struct {
	Port *Port `kiper_value:"name:port;help:port of server;default:8080"`
}
```

## Examples

[basic_usage.go](https://github.com/leeif/kiper/blob/master/_example/basic_usage.go)

[kiper_value.go](https://github.com/leeif/kiper/blob/master/_example/kiper_value.go)

## Limitation
* Config file settings have a higher priority when merge with flags
* Only support int and string array(slice) type now
* Not support sub command yet
