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
