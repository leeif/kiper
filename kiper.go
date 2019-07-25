package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

type KiperConfig interface {
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
		k.viper.SetConfigFile(path)
	}
	if err := k.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			// Config file was found but another error was produced
			return err
		}
	}
	return nil
}

func (k *Kiper) flags(config KiperConfig, kcName string) error {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		kc := field.Tag.Get("kiper_config")
		if kc != "" {
			if err := k.flags(value.Interface().(KiperConfig), kc); err != nil {
				return err
			}
			continue
		}

		kv := field.Tag.Get("kiper_value")
		if kv != "" {
			m := make(map[string]string)
			for _, k := range strings.Split(kv, ";") {
				keyPair := strings.Split(k, ":")
				if len(keyPair) < 2 {
					continue
				}
				m[keyPair[0]] = keyPair[1]
			}
			name, ok := m["name"]
			if !ok {
				continue
			}
			flag := ""
			if kcName != "" {
				flag = kcName + "." + name
			} else {
				flag = name
			}

			deflt, ok := m["default"]
			if !ok {
				deflt = ""
			}

			hp, ok := m["help"]
			if !ok {
				hp = ""
			}

			if field.Type.Kind() == reflect.Ptr {
				switch field.Type.Elem().Kind() {
				case reflect.String:
					s := k.kingpin.Flag(flag, hp).Default(deflt).String()
					v.Field(i).Set(reflect.ValueOf(s))
				case reflect.Int:
					s := k.kingpin.Flag(flag, hp).Default(deflt).Int()
					v.Field(i).Set(reflect.ValueOf(s))
				case reflect.Bool:
					s := k.kingpin.Flag(flag, hp).Default(deflt).Bool()
					v.Field(i).Set(reflect.ValueOf(s))
				case reflect.Struct:
					if v.Field(i).MethodByName("Set").IsValid() && v.Field(i).MethodByName("String").IsValid() {
						k.kingpin.Flag(flag, hp).Default(deflt).SetValue(v.Field(i).Interface().(KiperValue))
					}
				}
			} else if field.Type.Kind() == reflect.Struct {
				if v.Field(i).MethodByName("Set").IsValid() && v.Field(i).MethodByName("String").IsValid() {
					k.kingpin.Flag(flag, hp).Default(deflt).SetValue(v.Field(i).Interface().(KiperValue))
				}
			}
		}
	}

	return nil
}

func (k *Kiper) Merge(config KiperConfig) {
	k.merge(config, k.viper.AllSettings())
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
		if _, ok := m[kcName]; kcName != "" && ok {
			kc := v.Field(i).Interface().(KiperConfig)
			if reflect.TypeOf(m[kcName]).Kind() == reflect.Map {
				k.merge(kc, m[kcName].(map[string]interface{}))
			}
		}
		kvName := field.Tag.Get("kiper_value")
		if kvName != "" {
			d := make(map[string]string)
			for _, k := range strings.Split(kvName, ";") {
				keyPair := strings.Split(k, ":")
				if len(keyPair) < 2 {
					continue
				}
				d[keyPair[0]] = keyPair[1]
			}
			name, ok := d["name"]
			if !ok {
				continue
			}
			_, ok = m[name]
			if !ok {
				continue
			}

			if t.Field(i).Type.Kind() == reflect.Ptr {
				vField := v.Field(i).Elem()
				switch vField.Kind() {
				case reflect.String:
					vField.SetString(m[name].(string))
				case reflect.Int:
					vField.SetInt(m[name].(int64))
				case reflect.Bool:
					vField.SetBool(m[name].(bool))
				case reflect.Struct:
					if v.Field(i).MethodByName("Set").IsValid() && v.Field(i).MethodByName("String").IsValid() {
						v.Field(i).Interface().(KiperValue).Set(m[name].(string))
					}
				}
			} else if t.Field(i).Type.Kind() == reflect.Struct {
				if v.Field(i).MethodByName("Set").IsValid() && v.Field(i).MethodByName("String").IsValid() {
						v.Field(i).Interface().(KiperValue).Set(m[name].(string))
				}
			}
		}
	}
}

func NewKiper(config KiperConfig, name, help string) (*Kiper, error) {
	kiper := &Kiper{}
	kiper.viper = viper.New()
	kiper.kingpin = kingpin.New(name, help)

	if err := kiper.flags(config, ""); err != nil {
		return nil, err
	}

	return kiper, nil
}

type TestConfig struct {
	Address *Address `kiper_value:"name:address;help:address of server;default:127.0.0.1"`
	Test    *string  `kiper_value:"name:test;default:test"`
	Another Another  `kiper_config:"another"`
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

func main() {
	tc := &TestConfig{
		Address: &Address{},
		Another: Another{
			Address: &Address{},
		},
	}

	kiper, err := NewKiper(tc, "", "")
	if err != nil {
		fmt.Println(err)
	}

	kiper.SetCommandLineArguments(os.Args[1:])
	kiper.SetConfigFilePath("./config.json")

	err = kiper.Parse()
	if err != nil {
		fmt.Println(err)
	}
	kiper.Merge(tc)
	fmt.Println(tc.Address.String())
}
