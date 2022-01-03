package types

import (
	"github.com/duanchi/min/types/config"
	"reflect"
)

type Config struct {
	Env         string             `yaml:"env" default:"development"`
	Db          config.Db          `yaml:"db"`
	Application config.Application `yaml:"application"`
	Rpc         config.Rpc         `yaml:"rpc"`
	Feign       config.Feign       `yaml:"feign"`
	Log         config.Log         `yaml:"log"`
	Cache       config.Cache       `yaml:"cache"`
	Task        config.Task        `yaml:"task"`
	Discovery   config.Discovery   `yaml:"discovery"`
	BeanParsers interface{}        `yaml:"-"`
	Beans       struct{}           `yaml:"-"`
}

func (this *Config) GetName() (name string) {
	return "Config"
}
func (this *Config) SetName(name string) {
	return
}

type BeanParser struct {
}

func (parser BeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type) {}
