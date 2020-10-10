package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 解析yaml配置文件

type YamlConfig struct {
	//ModeDisk    []ModeDiskItem `yaml:"modeDisk"`
	ModeYum     []YumMode     `yaml:"modeYum"`
	ModeMysql   []MysqlMode   `yaml:"modeMysql"`
	ModeScripts []ScriptsMode `yaml:"modeScripts"`
}

//type ModeDiskItem struct {
//	Action string   `yaml:"action"`
//	Name   []string `yaml:"name"`
//}

type YumMode struct {
	Action string   `yaml:"action"`
	Name   []string `yaml:"name"`
}

type MysqlMode struct {
	Action string `yaml:"action"`
	Role   string `yaml:"role"`
	Host   string `yaml:"host"`
}

type ScriptsMode struct {
	Action string   `yaml:"action"`
	Env    string   `yaml:"env"`
	Name   []string `yaml:"name"`
}

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
