package kiper

import "reflect"

type KiperConfig interface {
	Name() string
}

type KiperValue interface {
	Set(string) error
	String() string
	Name() string
}

type Kiper struct {
	configFilePath []string
}

func (k *Kiper) SetConfigFilePath(path string) {
	k.configFilePath = append(k.configFilePath, path)
}

func (k *Kiper) Parse(config KiperConfig) {
}

func (k *Kiper) merge(config KiperConfig, map[string]interface{}) {
  t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
}
