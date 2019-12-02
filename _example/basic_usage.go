package main

import (
	"fmt"
	"os"

	"github.com/leeif/kiper"
)

type KiperTester struct {
	StringVar string    `kiper_value:"name:string_var;help:"string var";default:test string"`
	IntVar    int       `kiper_value:"name:int_var;default:10"`
	BoolVar   bool      `kiper_value:"name:bool_var;default:true"`
	ArrayVar  []string  `kiper_value:"name:array_var;default:test1, test2"`
	StructVar *StructVar `kiper_value:"name:struct_var;default:I am a struct"`
}

type StructVar struct {
	data string
}

func (structVar *StructVar) Set(data string) error {

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
	// set your own config file path flag with default value
	k.SetConfigFileFlag("config", "config file", "./config.json")
	k.Kingpin.HelpFlag.Short('h')

	// parse command line and config file
	if err := k.Parse(kt, os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(kt.StringVar)
	fmt.Println(kt.IntVar)
	fmt.Println(kt.BoolVar)
	fmt.Println(kt.StructVar.data)
	fmt.Println(kt.ArrayVar)
}
