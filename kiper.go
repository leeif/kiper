package kiper

import (
	"errors"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

type KiperValue interface {
	Set(string) error
	String() string
}

type Kiper struct {
	viper   *viper.Viper
	kingpin *kingpin.Application
}

func (k *Kiper) GetViperInstance() *viper.Viper {
	return k.viper
}

func (k *Kiper) GetKingpinInstance() *kingpin.Application {
	return k.kingpin
}

func (k *Kiper) ParseCommandLine(config interface{}, args []string) error {
	startKiperConfig := ""
	k.flags(config, startKiperConfig)
	// parse command line flags
	if _, err := k.kingpin.Parse(args); err != nil {
		return err
	}

	return nil
}

func (k *Kiper) ParseConfigFile(path ...string) error {
	for _, path := range path {
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

func (k *Kiper) flags(config interface{}, kcName string) error {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return errors.New("Kiper Config " + t.Name() + " is not Struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tags := k.parseTag(field.Tag.Get("kiper_config"))
		if name, ok := tags["name"]; ok && name != "" {
			if err := k.flags(value.Interface(), name); err != nil {
				return err
			}
			continue
		}

		tags = k.parseTag(field.Tag.Get("kiper_value"))

		if name, ok := tags["name"]; !ok || name == "" {
			continue
		}
		hp, deflt := tags["help"], tags["default"]

		flag := ""
		if kcName == "" {
			flag = tags["name"]
		} else {
			flag = kcName + "." + tags["name"]
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

	return nil
}

func (k *Kiper) parseTag(tag string) map[string]string {
	m := make(map[string]string)
	for _, k := range strings.Split(tag, ";") {
		keyPair := strings.Split(k, ":")
		if len(keyPair) < 2 {
			continue
		}
		m[keyPair[0]] = keyPair[1]
	}
	return m
}

func (k *Kiper) MergeConfigFile(config interface{}) error {

	// get all config file setting
	if err := k.merge(config, k.viper.AllSettings()); err != nil {
		return err
	}
	return nil
}

func (k *Kiper) merge(config interface{}, m map[string]interface{}) error {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return errors.New("Config is not Struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tags := k.parseTag(field.Tag.Get("kiper_config"))
		if name, ok := tags["name"]; ok && name != "" {
			if v, ok := m[name]; !ok || reflect.TypeOf(v).Kind() != reflect.Map {
				continue
			}
			if field.Type.Kind() != reflect.Struct && field.Type.Kind() != reflect.Ptr {
				continue
			}
			if err := k.merge(value.Interface(), m[name].(map[string]interface{})); err != nil {
				return err
			}
			continue
		}
		tags = k.parseTag(field.Tag.Get("kiper_value"))
		if name, ok := tags["name"]; ok && name != "" {
			if _, ok = m[name]; !ok {
				continue
			}
			if t.Field(i).Type.Kind() == reflect.Ptr {
				vField := v.Field(i).Elem()
				switch vField.Kind() {
				case reflect.String:
					vField.SetString(m[name].(string))
				case reflect.Int:
					vField.SetInt(int64(m[name].(float64)))
				case reflect.Bool:
					vField.SetBool(m[name].(bool))
				case reflect.Struct:
					if v.Field(i).MethodByName("Set").IsValid() && v.Field(i).MethodByName("String").IsValid() {
						if err := v.Field(i).Interface().(KiperValue).Set(m[name].(string)); err != nil {
							return err
						}
					}
				}
			} else if t.Field(i).Type.Kind() == reflect.Struct {
				if v.Field(i).MethodByName("Set").IsValid() && v.Field(i).MethodByName("String").IsValid() {
					if err := v.Field(i).Interface().(KiperValue).Set(m[name].(string)); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func NewKiper(name, help string) *Kiper {
	kiper := &Kiper{}
	kiper.viper = viper.New()
	kiper.kingpin = kingpin.New(name, help)

	return kiper
}
