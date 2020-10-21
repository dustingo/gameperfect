package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

// hostname

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
	Host   string      `yaml:"host"`
}

type ScriptsMode struct {
	Action string   `yaml:"action"`
	Env    string   `yaml:"env"`
	Name   []string `yaml:"name"`
	Host   string   `yaml:"host"`
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

// config.yaml中执行yum操作
func DoYumAction(yamlConfig *YamlConfig) {
	if yamlConfig.ModeYum != nil {
		fmt.Printf(CSI+Green+"%s"+End+"\n", "Found Mod: Yum")
		for _, act := range yamlConfig.ModeYum {
			if act.Action == "install" {
				for i := 0; i < len(act.Name); i++ {
					if YumCheck(act.Name[i]) {
						fmt.Println(CSI + Blue + "[info] " + act.Name[i] + " already installed" + End)
					} else {
						fmt.Println(CSI + Blue + ">Do Install:" + act.Name[i] + End)
						time.Sleep(time.Second * 1)
						err := Execute("yum", "-y", "install", act.Name[i])
						if err != nil {
							fmt.Printf(CSI+Red+"[Error] "+"%s"+"\n", err)
						}
					}

				}

			}
		}
	}
}

// config.yaml 执行scripts操作
func DoScriptAction(yamlConfig *YamlConfig) {
	host, _ := os.Hostname()
	if yamlConfig.ModeScripts != nil {
		fmt.Printf(CSI+Green+"%s"+End+"\n", "Found Mod: Scripts")
		for _, act := range yamlConfig.ModeScripts {
			if act.Action == "run" {
				if host == act.Host || act.Host == "all" || act.Host == "" {
					if act.Env == "sh" || act.Env == "bash" {
						for i := 0; i < len(act.Name); i++ {
							fmt.Println(CSI + Blue + ">Run " + act.Name[i] + End)
							err := Execute(act.Env, "-c", act.Name[i])
							if err != nil {
								fmt.Println(CSI + Red + "[Error] " + err.Error() + End)
								//return
							}
						}
					} else if act.Env == "python" || act.Env == "Python" {
						for i := 0; i < len(act.Name); i++ {
							fmt.Println(CSI + Blue + ">Run " + act.Name[i] + End)
							err := Execute(act.Env, act.Name[i])
							if err != nil {
								fmt.Println(CSI + Red + "[Error] " + err.Error() + End)
								//return
							}

						}
					}
				} else {
					fmt.Println(CSI + Red + "[Info] " + "host not " + act.Host + " and " + act.Env + " will not run" + End)
				}
			} else {
				fmt.Println(CSI + Red + "[Error] " + act.Action + " can't run" + End)
			}
		}
	}

}

// 执行config.yaml中 dir
func DoDirAction(yamlConfig *YamlConfig) {
	host, _ := os.Hostname()
	if yamlConfig.ModeDir != nil {
		fmt.Printf(CSI+Green+"%s"+End+"\n", "Found Mod: Dir")
		for _, act := range yamlConfig.ModeDir {
			if act.Action == "chown" {
				if act.Host == host || act.Host == "all" || act.Host == "" {
					for i := 0; i < len(act.Path); i++ {
						fmt.Println(CSI + Blue + ">Do chown:" + act.Path[i] + End)
						time.Sleep(time.Second)
						if PathExists(act.Path[i]) {
							err := Execute(act.Action, act.Para, act.Path[i])
							if err != nil {
								fmt.Printf(CSI+Red+"[Error] "+"%s"+"\n", err)
							}
						} else {
							Mkdir(act.Path[i], 0755)
						}

					}
				} else {
					fmt.Println(CSI + Green + "[Info]" + "hostname not " + act.Host + " " + act.Action + " will not run " + End)
				}
			} else if act.Action == "chmod" {
				if act.Host == host || act.Host == "all" || act.Host == "" {
					for i := 0; i < len(act.Path); i++ {
						fmt.Println(CSI + Blue + ">Do chmod:" + act.Path[i] + End)
						time.Sleep(time.Second)
						if PathExists(act.Path[i]) {
							//fmt.Println(act.Perm)
							//err := config.Execute(act.Action, act.Perm, act.Path[i])
							err := os.Chmod(act.Path[i], act.Perm)
							if err != nil {
								fmt.Printf(CSI+Red+"%s"+"\n", err)
							}
						} else {
							Mkdir(act.Path[i], act.Perm)
						}

					}
				} else {
					fmt.Println(CSI + Green + "[Info]" + "hostname not " + act.Host + " " + act.Action + " will not  run " + End)
				}
			}
		}
	}
}
