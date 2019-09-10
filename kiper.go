package kiper

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

// DefaultConfigFile is default config file path
const (
	DefaultConfigFile = ".kiper"
)

type KiperValue interface {
	Set(string) error
	String() string
}

type kiperValue struct {
	data string
}

func (kv *kiperValue) Set(data string) error {
	kv.data = data
	return nil
}

func (kv *kiperValue) String() string {
	return kv.data
}

type Kiper struct {
	Viper      *viper.Viper
	Kingpin    *kingpin.Application
	configFile *string
	// kiper config map map struct
	kpMap map[string]interface{}

	// viper config map
	vpMap map[string]interface{}
}

func (k *Kiper) Parse(config interface{}, args []string) error {
	var err error
	startKiperConfig := ""
	k.kpMap, err = k.parseFlags(config, startKiperConfig)
	if err != nil {
		return err
	}

	// parse command line flags
	if _, err := k.Kingpin.Parse(args); err != nil {
		return err
	}

	fmt.Println(*k.configFile)
	if k.configFile != nil && *k.configFile != "" {
		k.vpMap, err = k.parseConfigFile(*k.configFile)
	} else {
		k.vpMap, err = k.parseConfigFile(DefaultConfigFile)
	}
	if err != nil {
		return err
	}

	err = k.merge(config, k.kpMap, k.vpMap)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kiper) SetConfigFileFlag(flag string, description string, value string) {
	if value != "" {
		k.configFile = k.Kingpin.Flag(flag, description).Default(value).String()
		return
	}
	k.configFile = k.Kingpin.Flag(flag, description).String()
}

func (k *Kiper) parseConfigFile(path string) (map[string]interface{}, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// return implicitly if file is not exists
		return nil, nil
	}

	k.Viper.SetConfigFile(path)
	if err := k.Viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return k.Viper.AllSettings(), nil
}

func (k *Kiper) parseFlags(config interface{}, kcName string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("Kiper Config " + t.Name() + " is not Struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tags := k.parseTag(field.Tag.Get("kiper_config"))
		if name, ok := tags["name"]; ok && name != "" {
			m, err := k.parseFlags(value.Interface(), name)
			if err != nil {
				return nil, err
			}
			res[name] = m
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

		switch field.Type.Kind() {
		case reflect.String:
			res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).String()
		case reflect.Int:
			res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).Int()
		case reflect.Bool:
			res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).Bool()
		case reflect.Struct:
			res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).String()
		case reflect.Ptr:
			switch field.Type.Elem().Kind() {
			case reflect.String:
				res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).String()
			case reflect.Int:
				res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).Int()
			case reflect.Bool:
				res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).Bool()
			case reflect.Struct:
				res[tags["name"]] = k.Kingpin.Flag(flag, hp).Default(deflt).String()
			}
		}
	}

	return res, nil
}

func (k *Kiper) parseTag(tag string) map[string]string {
	m := make(map[string]string)
	for _, k := range strings.Split(tag, ";") {
		keyPair := strings.Split(k, ":")
		if len(keyPair) < 2 {
			continue
		}
		// rejoin the rest part of tag using `:`
		m[keyPair[0]] = strings.Join(keyPair[1:], ":")
	}
	return m
}

func (k *Kiper) merge(config interface{}, kpMap map[string]interface{}, vpMap map[string]interface{}) error {
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
			if field.Type.Kind() != reflect.Struct && field.Type.Kind() != reflect.Ptr {
				continue
			}
			var ok bool
			var vm map[string]interface{}
			var km map[string]interface{}
			vm, ok = vpMap[name].(map[string]interface{})
			if !ok {
				vm = nil
			}
			km, ok = kpMap[name].(map[string]interface{})
			if !ok {
				km = nil
			}
			if err := k.merge(value.Interface(), km, vm); err != nil {
				return err
			}
			continue
		}
		tags = k.parseTag(field.Tag.Get("kiper_value"))
		if name, ok := tags["name"]; ok && name != "" {
			switch field.Type.Kind() {
			case reflect.String:
				k.setStringValue(value, kpMap[name], vpMap[name])
			case reflect.Int:
				k.setIntValue(value, kpMap[name], vpMap[name])
			case reflect.Bool:
				k.setBoolValue(value, kpMap[name], vpMap[name])
			case reflect.Struct:
				k.setKiperValue(value, kpMap[name], vpMap[name])
			case reflect.Ptr:
				err := k.setPointerValue(value, field.Type, kpMap[name], vpMap[name])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (k *Kiper) setStringValue(value reflect.Value, flag interface{}, cfg interface{}) {
	if flag == nil && cfg == nil {
		return
	}
	fv, ok1 := flag.(*string)
	cv, ok2 := cfg.(string)
	if ok2 {
		value.SetString(cv)
		return
	}
	if ok1 {
		value.SetString(*fv)
		return
	}
}

func (k *Kiper) setIntValue(value reflect.Value, flag interface{}, cfg interface{}) {
	if flag == nil && cfg == nil {
		return
	}
	fv, ok1 := flag.(*int)
	cv, ok2 := cfg.(float64)
	if ok2 {
		value.SetInt(int64(cv))
		return
	}
	if ok1 {
		value.SetInt(int64(*fv))
		return
	}
}

func (k *Kiper) setBoolValue(value reflect.Value, flag interface{}, cfg interface{}) {
	if flag == nil && cfg == nil {
		return
	}
	fv, ok1 := flag.(*bool)
	cv, ok2 := cfg.(bool)
	if ok2 {
		value.SetBool(cv)
		return
	}
	if ok1 {
		value.SetBool(*fv)
		return
	}
}

func (k *Kiper) setKiperValue(value reflect.Value, flag interface{}, cfg interface{}) error {
	if flag == nil && cfg == nil {
		return nil
	}

	kv, ok := value.Interface().(KiperValue)
	if !ok {
		return nil
	}

	fv, ok1 := flag.(*string)
	cv, ok2 := cfg.(string)
	var err error
	if ok2 {
		err = kv.Set(cv)
		if err != nil {
			return err
		}
		return nil
	}
	if ok1 {
		err = kv.Set(*fv)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (k *Kiper) setPointerValue(value reflect.Value, t reflect.Type, flag interface{}, cfg interface{}) error {
	switch t.Elem().Kind() {
	case reflect.String:
		k.setPointerString(value, flag, cfg)
	case reflect.Int:
		k.setPointerInt(value, flag, cfg)
	case reflect.Bool:
		k.setPointerBool(value, flag, cfg)
	case reflect.Struct:
		err := k.setKiperValue(value, flag, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kiper) setPointerString(value reflect.Value, flag interface{}, cfg interface{}) {
	if flag == nil && cfg == nil {
		return
	}
	fv, ok1 := flag.(*string)
	cv, ok2 := cfg.(string)
	if ok2 {
		value.Set(reflect.ValueOf(&cv))
		return
	}
	if ok1 {
		value.Set(reflect.ValueOf(fv))
		return
	}
}

func (k *Kiper) setPointerInt(value reflect.Value, flag interface{}, cfg interface{}) {
	if flag == nil && cfg == nil {
		return
	}
	fv, ok1 := flag.(*int)
	cv, ok2 := cfg.(float64)
	if ok2 {
		tmp := int(cv)
		value.Set(reflect.ValueOf(&tmp))
		return
	}
	if ok1 {
		value.Set(reflect.ValueOf(fv))
		return
	}
}

func (k *Kiper) setPointerBool(value reflect.Value, flag interface{}, cfg interface{}) {
	if flag == nil && cfg == nil {
		return
	}
	fv, ok1 := flag.(*bool)
	cv, ok2 := cfg.(bool)
	if ok2 {
		value.Set(reflect.ValueOf(&cv))
		return
	}
	if ok1 {
		value.Set(reflect.ValueOf(fv))
		return
	}
}

func NewKiper(name, help string) *Kiper {
	kiper := &Kiper{}
	kiper.Viper = viper.New()
	kiper.Kingpin = kingpin.New(name, help)

	return kiper
}
