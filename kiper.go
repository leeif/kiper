package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

type KiperConfig interface {
	KCName() string
}

type KiperValue interface {
	Set(string) error
	String() string
}

type Kiper struct {
	viper          *viper.Viper
	kingpin        *kingpin.Application
	configFilePath []string
	arguments      []string
}

func (k *Kiper) SetConfigFilePath(path string) {
	k.configFilePath = append(k.configFilePath, path)
}

func (k *Kiper) SetCommandLineArguments(args []string) {
	fmt.Println(args)
	k.arguments = args
}

func (k *Kiper) GetViperInstance() *viper.Viper {
	return k.viper
}

func (k *Kiper) GetKingpinInstance() *kingpin.Application {
	return k.kingpin
}

func (k *Kiper) Parse() error {
	if err := k.configFile(); err != nil {
		return err
	}

	k.kingpin.Parse(k.arguments)

	return nil
}

func (k *Kiper) configFile() error {
	for _, path := range k.configFilePath {
		viper.SetConfigFile(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			// Config file was found but another error was produced
			return err
		}
	}
	return nil
}

func (k *Kiper) flags(config KiperConfig) error {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if value. != "" {
			k.flags(v.Field(i).Interface().(KiperConfig))
		}

		kvName := field.Tag.Get("kiper_value")
		if kvName != "" {
			if v.Field(i).MethodByName("Set").IsValid() && v.Field(i).MethodByName("String").IsValid() {
				k.kingpin.Flag(config.Name()+"."+kvName, "").SetValue(v.Field(i).Interface().(KiperValue))
				continue
			}
			fmt.Println(kvName)
			s := k.kingpin.Flag(config.Name()+"."+kvName, "").String()
			v.Field(i).Set(reflect.ValueOf(s))
		}
	}

	return nil
}

func (k *Kiper) merge(config KiperConfig, m map[string]interface{}) {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		kcName := field.Tag.Get("kiper_config")
		if kcName != "" && field.Type.Kind() == reflect.Struct {
			kc := v.Field(i).Interface().(KiperConfig)
			if reflect.TypeOf(m[kcName]).Kind() == reflect.Map {
				k.merge(kc, m[kcName].(map[string]interface{}))
			}
		}
		kvName := field.Tag.Get("kiper_value")
		if kvName != "" {
			switch t.Field(i).Type.Kind() {
			case reflect.String:
				v.Field(i).SetString(m[kvName].(string))
			case reflect.Int:
				v.Field(i).SetInt(m[kvName].(int64))
			case reflect.Struct:
				v.Field(i).Interface().(KiperValue).Set(m[kvName].(string))
			}
		}
		fmt.Println(v.MethodByName("asdasd").IsValid())
	}
}

func NewConfig(config KiperConfig, name, help string) (*Kiper, error) {
	kiper := &Kiper{}
	kiper.viper = viper.New()
	kiper.kingpin = kingpin.New(name, help)

	if err := kiper.flags(config); err != nil {
		return nil, err
	}

	return kiper, nil
}

type TestConfig struct {
	Address *Address `kiper_value:"name:address;help:address of server;default:127.0.0.1"`
	Test    *string  `kiper_value:"name:test"`
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

func (t *TestConfig) Name() string {
	return "test_config"
}

func main() {
	tc := &TestConfig{
		Address: &Address{},
	}

	kiper, err := NewConfig(tc, "", "")
	if err != nil {
		fmt.Println(err)
	}

	kiper.SetCommandLineArguments(os.Args[1:])

	err = kiper.Parse()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tc.Test)
}
