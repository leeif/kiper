package main

import (
	"testing"

	"github.com/leeif/kiper"
)

type TestConfig struct {
}

func (t *TestConfig) Name() string {
	return "test_config"
}

func TestParse(t *testing.T) {
	kiper := kiper.New("", "")
	tc := &TestConfig{}
	err := kiper.Parse(tc)
	if err != nil {
		t.Error(err)
	}
}
