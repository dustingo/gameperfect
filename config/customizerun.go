package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// 解析yaml配置文件

type YamlConfig struct {
	ModeYum     []YumMode     `yaml:"modeYum"`
	ModeDir     []DirMode     `yaml:"modeDir"`
	ModeScripts []ScriptsMode `yaml:"modeScripts"`
}

type YumMode struct {
	Action string   `yaml:"action"`
	Name   []string `yaml:"name"`
}

type DirMode struct {
	Action string      `yaml:"action"`
	Para   string      `yaml:"para"`
	Path   []string    `yaml:"path"`
	Perm   os.FileMode `yaml:"perm"`
}

type ScriptsMode struct {
	Action string   `yaml:"action"`
	Env    string   `yaml:"env"`
	Name   []string `yaml:"name"`
}

// ParseYaml 解析config.yaml
func ParseYaml() *YamlConfig {
	ymlConfig := new(YamlConfig)
	yamlByte, err := ioutil.ReadFile("./config/config.yaml")
	if len(yamlByte) == 0 {
		fmt.Println("config is null")
		return nil
	}
	if err != nil {
		fmt.Println("")
		return nil
	}
	err = yaml.Unmarshal(yamlByte, ymlConfig)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ymlConfig
}
