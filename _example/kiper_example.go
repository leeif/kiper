package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/leeif/kiper"
)

type Server struct {
	Address *Address `kiper_value:"name:address;default:127.0.0.1"`
	Port    *Port    `kiper_value:"name:port;default:3306"`
}

type Address struct {
	s string
}

func (address *Address) Set(s string) error {
	if s == "" {
		return errors.New("address can't be empty")
	}
	address.s = s
	return nil
}

func (address *Address) String() string {
	return address.s
}

type Port struct {
	p string
}

func (port *Port) Set(p string) error {
	if _, err := strconv.Atoi(p); err != nil {
		return errors.New("not a valid port value")
	}
	port.p = p
	return nil
}

func (port *Port) String() string {
	return port.p
}

type Config struct {
	ConfigFile *string `kiper_value:"name:config_file;default:./config.json"`
	ID         *int    `kiper_value:"name:id;default:1"`
	Server     Server  `kiper_config:"name:server"`
}

func main() {
	// initialize config struct
	c := &Config{
		Server: Server{
			Address: &Address{},
			Port:    &Port{},
		},
	}

	// new kiper
	k := kiper.NewKiper("example", "example of kiper")

	// parse command line and config file
	if err := k.ParseCommandLine(c, os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(c.Server.Address)
	fmt.Println(*c.ID)

	// set config file path for viper
	k.ParseConfigFile("./config.json")

	// merge with config file
	if err := k.MergeConfigFile(c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(c.Server.Address)
	fmt.Println(*c.ID)
}
