package kiper_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/leeif/kiper"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Address    *Address `kiper_value:"name:address;help:address of server;default:127.0.0.1"`
	TestString string   `kiper_value:"name:test_string;default:test_string"`
	TestInt    int      `kiper_value:"name:test_int;required"`
	TestInt64  int64    `kiper_value:"name:test_int64;default:64"`
	TestBool   bool     `kiper_value:"name:test_bool;default:false"`
	Another    Another  `kiper_config:"name:another"`
}

type Another struct {
	Address *Address `kiper_value:"name:address;help:address of server;default:127.0.0.1"`
}

type Address struct {
	s string
}

func (a *Address) Set(s string) error {
	a.s = s
	return nil
}

func (a *Address) String() string {
	return a.s
}

func writeConfigFile(path string, s interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func deleteConfigFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func TestKiperConfig(t *testing.T) {
	tc := &TestConfig{
		Address: &Address{},
		Another: Another{
			Address: &Address{},
		},
	}

	// command line flags
	args := []string{
		"--config", "./config.json",
		"--address", "10.0.0.1",
		"--test_string", "not test",
		"--test_int", "2",
		"--test_int64", "80",
		"--another.address", "192.0.0.1",
	}

	// config file
	writeConfigFile("./config.json", struct {
		Address    string `json:"address"`
		TestString string `json:"test_string"`
		Another    struct {
			Address string `json:"address"`
		} `json:"another"`
	}{
		Address:    "test1",
		TestString: "test2",
		Another: struct {
			Address string `json:"address"`
		}{
			Address: "test3",
		},
	})
	defer deleteConfigFile("./config.json")

	kiper := kiper.NewKiper("test", "kiper test")
	kiper.SetConfigFileFlag("config", "config file flag", "./config.json")
	err := kiper.Parse(tc, args)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, tc.TestInt, 2, "test should be test2")
	assert.Equal(t, tc.TestBool, false, "test should be test2")

	assert.Equal(t, tc.TestInt64, int64(80), "test should be 80")

	assert.Equal(t, tc.Address.String(), "test1", "address should be test1 after merge")
	assert.Equal(t, "test2", tc.TestString, "test should be test2 after merge")
	assert.Equal(t, tc.Another.Address.String(), "test3", "another.address should be test3 after merge")
}
