package lib

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

type icon struct {
	Width uint
	Height uint
	Radius uint
	Type string
}

type configData struct {
	Prefix string
	Icons []icon
}

type ConfigProvider struct {
	ConfigData configData
}

func (cp *ConfigProvider) Initialize(configFile string) error {
	cp.ConfigData = configData{}

	if data, err := ioutil.ReadFile(configFile); err != nil {
		return err
	} else {
		yaml.Unmarshal(data, &cp.ConfigData)
		fmt.Println(cp.ConfigData)
	}

	return nil
}
