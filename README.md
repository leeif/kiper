# kiper
A small wrapper of [kingpin](https://github.com/alecthomas/kingpin) + [viper](https://github.com/spf13/viper.git). I like the way viper handles the config files however personally I am a fan of kinping because it's fluent-style and it can validate the command line flags easily. This library is a combination of this two projects to handle the flags and config files at the same time.

# Feature

* Merge flag and config file settings automatically.
* Validate the settings for both flag and config file settings.

# Usage

```
$ go get github.com/leeif/kiper
```

## Simple Example

```
package main

import (
	"fmt"
	"os"

	"github.com/leeif/kiper"
)

type Config struct {
	Name    *string `kiper_value:"name:name;required"`
	Gender  *int    `kiper_value:"name:gender;default:1"`
	FanOfGo *bool   `kiper_value:"name:fan_of_go;default:false"`
}

func main() {
	// initialize config struct
	c := &Config{}

	// new kiper
	k := kiper.NewKiper("example", "example of kiper")

	// parse command line and config file
	if err := k.Parse(c, os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("name: %s\n", *c.Name)
	fmt.Printf("gender: %d\n", *c.Gender)
	fmt.Printf("fan of go: %v\n", *c.FanOfGo)
}
```

```
go run . --name=leeif
name: leeif
gender: 1
fan of go: false
```

### Example of merging flag and config file settings

[kiper_example.go](https://github.com/leeif/kiper/blob/master/_example/kiper_example.go)

```
$ dep ensure
$ cd _example
$ go run kiper_example.go --help
```

## kiper_value

kiper_value now support the types of int, string, bool and their pointer types.

kiper_value can also be a struct which implement the interface for validations.

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
|  required      | set the flag to required, default value will be ignored due to the kingpin feature |

## Limitation
* Config file settings have a higher priority when merge with flags
* Not support sub command yet
