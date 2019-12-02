package main

import (
	"fmt"
	"os"

	"github.com/leeif/kiper"
)

type KiperTester struct {
	StructVar *StructVar `kiper_value:"name:struct_var"`
}

type StructVar struct {
	data string
}

func (structVar *StructVar) Set(data string) error {
	if data == "" {
		return fmt.Errorf("data can not be empty")
	}
	structVar.data = data
	return nil
}

func (structVar *StructVar) String() string {
	return structVar.data
}

func main() {
	// initialize config struct
	kt := &KiperTester{
		StructVar: &StructVar{},
	}

	// new kiper
	k := kiper.NewKiper("example", "example of kiper")
	k.Kingpin.HelpFlag.Short('h')

	// parse command line and config file
	if err := k.Parse(kt, os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(kt.StructVar.data)
}
