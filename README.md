# kiper
A small wrapper of [kingpin](https://github.com/alecthomas/kingpin) + [viper](https://github.com/spf13/viper.git). I like the way viper handles the config files however I am a fan of kinping because it's fluent-style and it can validate the command line flags easily. This library is just a combination of the features of this two projects.

# Feature

* Merge command line flags and config file settings.
* Validate command line flags and config file bothly.

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
	Name    *string `kiper_value:"name:name;default:leeif"`
	Gender  *int    `kiper_value:"name:gender;default:1"`
	FanOfGo *bool   `kiper_value:"name:fan_of_go;default:false"`
}

func main() {
	// initialize config struct
	c := &Config{}

	// new kiper
	k := kiper.NewKiper("example", "example of kiper")

	// parse command line and config file
	if err := k.ParseCommandLine(c, os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("name: %s\n", *c.Name)
	fmt.Printf("gender: %d\n", *c.Gender)
	fmt.Printf("fan of go: %v\n", *c.FanOfGo)
}
```

output

```
$ go run ./ --help
usage: example [<flags>]

example of kiper

Flags:
  --help          Show context-sensitive help (also try --help-long and --help-man).
  --name="leeif"
  --gender=1
  --fan_of_go

$ go run ./
name: leeif
gender: 1
fan of go: false

// Noop, I am a fan
$ go run ./ --fan_of_go --fan_of_go
name: leeif
gender: 1
fan of go: true
```

### Example of combination with command line flags and config file

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

## Limitation

* All the primitive types should be pointer types.
* Not support sub command yet
