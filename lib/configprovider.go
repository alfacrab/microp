package lib

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type configData struct {
	Sets []IconSet
}

func (data *configData) Length() int {
	result := 0

	for _, set := range data.Sets {
		result += len(set.Icons)
	}

	return result
}

type IconSet struct {
	Prefix string
	Icons []IconConfig
}

type IconConfig struct {
	Width uint
	Height uint
	Radius uint
	Type string
  Name string
}

type ConfigProvider struct {
	ConfigFile string
	ConfigData configData
}

func (cp *ConfigProvider) Initialize(configFile string) error {
	cp.ConfigData = configData{}

	if data, err := ioutil.ReadFile(configFile); err != nil {
		return err
	} else {
		err = yaml.Unmarshal(data, &cp.ConfigData)
    if err != nil {
      return err
    }

		cp.ConfigFile = configFile
	}

	return nil
}
