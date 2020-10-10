package main

import (
	"bytes"
	"fmt"
	"gameperfect/config"
	"time"

	//ui "github.com/gizak/termui/v3"
	//"github.com/gizak/termui/v3/widgets"
	"gopkg.in/alecthomas/kingpin.v2"
	//"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	doYaml    = kingpin.Flag("yaml", "only do yaml action，default is only do basic check[sys,cpu,disk ...]").Bool()
	doAll     = kingpin.Flag("all", "do basic check and yaml action").Bool()
	doHelp    = kingpin.Flag("usage", "The explanation of the config of this tool").Bool()
	doVersion = kingpin.Flag("version", "show the version of the tools").Bool()
)
var (
	version = "1.0.0"
)

func usage() {
	fmt.Println(">>目录：config 存放yaml配置文件")
	fmt.Println(">>目录:tools 存放自定义脚本")
	fmt.Println(">>config.yaml: ")
	fmt.Println("  >1.modeYum: action仅为'install'时，执行安装")
	fmt.Println("  >2.modeMysql: action为'install'时，根据role [master|slave],主机host选择安装")
	fmt.Println("  >3.modeScripts: action为'run'时，根据env [shell|python]和name [path + scripts_name]执行自定义脚本")
	fmt.Println(config.CSI + config.Red + "ATTENTION！注意：" + config.End)
	fmt.Println(config.CSI + config.Red + "the scripts should add 'x' permission" + config.End)
	fmt.Println(config.CSI + config.Red + "如果scripts是在win编辑的，在linux下记得使用dos转换" + config.End)
	fmt.Println(">>dns.yaml: ")
	fmt.Println("  >1.https的完整url，执行https GET")
	fmt.Println("  >2.域名，执行域名解析")
}

// 显示版本信息
func showVersion() {
	fmt.Println(version)
}

// 显示系统信息
func showSystem() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************SYS INFO********************")
	gameSystem := config.SystemInfo{}
	gameSystem.GetSystem()
	OsName := gameSystem.OSName
	OsRelease := gameSystem.OSRelease
	OsKernel := gameSystem.OSKernel
	OsArch := gameSystem.OSArch
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "OS Name: ", OsName)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "OS Release: ", OsRelease)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "OS Kernel: ", OsKernel)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "OS Arch: ", OsArch)
}

// 显示CPU信息
func showCpu() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************CPU INFO********************")
	newMap := make(map[string]int)
	gameCpu := config.CpuInfo{}
	gameCpu.GetCpu()
	cpuName := gameCpu.ModelName
	for i := 0; i < len(gameCpu.PhysicalID); i++ {
		newMap[gameCpu.PhysicalID[i]] = i
	}
	physicalId := len(newMap)
	cpuCores := gameCpu.CpuCores
	processor := len(gameCpu.Processor)
	cpuMhz := gameCpu.CpuMHz
	phyInt, _ := strconv.Atoi(cpuCores)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "CPU Name(CPU型号): ", cpuName)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%d"+"\n", "CPU Physical Number(物理CPU个数): ", physicalId)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "CPU Cores Number:(单颗CPU核心数): ", cpuCores)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%d"+"\n", "CPU Processor Number(CPU总逻辑个数): ", processor)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "CPU MHZ(CPU频率): ", cpuMhz)
	if physicalId*phyInt*2 == processor {
		hyperThreading := "YES"
		fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "HT（是否超线程）: ", hyperThreading)
	} else {
		hyperThreading := "NO"
		fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "HT（是否超线程）: ", hyperThreading)
	}

}

// 显示时区信息
func showTimeZone() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************TIME INFO********************")
	timeZone := config.Timezone{}
	timeZone.GetTimeZone()
	LocalTime := timeZone.LocalTime
	TimeZone := timeZone.Zone
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "LocalTime: ", LocalTime)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "TimeZone: ", TimeZone)

}

// 显示网卡信息 name,speed,ipaddress
func showNet() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************NET INFO********************")
	netFace := config.NetFaceInfo{}
	netFace.GetNetInfo()
	Names := netFace.Name
	Speeds := netFace.Speed
	Ipaddresses := netFace.Ipaddress
	for i := 0; i < len(Names); i++ {
		fmt.Printf("\033[1;36;40m%s\033[0m %s \033[1;36;40m%s\033[0m %s \033[1;36;40m%s\033[0m %s \n", "Name:", Names[i], "Speed:", Speeds[i], "Ipaddress:", Ipaddresses[i])
	}
}

// 显示硬盘信息
func showDisk() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************Disk INFO********************")
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"\n", "Filesystem              Type      Size  Used Avail Use% Mounted on")
	//fmt.Println("Filesystem              Type      Size  Used Avail Use% Mounted on")
	cmd := exec.Command("df", "-Th")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range strings.Split(out.String(), "\n") {
		if ok, _ := regexp.MatchString(config.IgnoreMountPoints, item); ok {
			fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"\n", item)

		}

	}
}

// 显示内存信息
func showMem() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************MEM INFO********************")
	memory := config.MemoryInfo{}
	memory.GetMemInfo()
	Total := memory.MemTotal
	Available := memory.MemAvailable
	Swapped := memory.MemSwap
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "MemTotal:", Total)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "MemAvailable:", Available)
	fmt.Printf(config.CSI+config.Blue+"%s"+config.End+"%s"+"\n", "SwapCached:", Swapped)
}

// 显示DNS信息
func showDns() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************DNS INFO********************")
	config.GetDnsInfo()
	config.ResolveDns()
}

// 显示NTP信息
func showNtp() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************NTP INFO********************")
	config.GetNtpInfo()
}

// 显示iptables信息
func showIptables() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************IPtables INFO********************")
	config.GetIptables()
}

// 显示ipmitool日志
func showIpmi() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************IPMI INFO********************")
	config.GetIpmiInfo()
}

// 执行config.yaml
func doYamlAction() {
	fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "********************Yaml Action********************")
	yamlConfig := config.ParseYaml()
	if yamlConfig.ModeYum != nil {
		fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "Found Mod: Yum")
		for _, act := range yamlConfig.ModeYum {
			if act.Action == "install" {
				for i := 0; i < len(act.Name); i++ {
					if config.YumCheck(act.Name[i]) {
						fmt.Println(config.CSI + config.Blue + act.Name[i] + " already installed" + config.End)
					} else {
						fmt.Println(config.CSI + config.Blue + ">Do Install:" + act.Name[i] + config.End)
						time.Sleep(time.Second * 1)
						err := config.Execute("yum", "-y", "install", act.Name[i])
						if err != nil {
							fmt.Printf(config.CSI+config.Red+"%s"+"\n", err)
						}
					}

				}

			}
		}
	}
	if yamlConfig.ModeMysql != nil {
		fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "Found Mod: MySQL")
	}
	if yamlConfig.ModeScripts != nil {
		fmt.Printf(config.CSI+config.Green+"%s"+config.End+"\n", "Found Mod: Scripts")
		for _, act := range yamlConfig.ModeScripts {
			if act.Action == "run" {
				if act.Env == "sh" || act.Env == "bash" {
					for i := 0; i < len(act.Name); i++ {
						fmt.Println(config.CSI + config.Blue + ">Run " + act.Name[i] + config.End)
						err := config.Execute(act.Env, "-c", act.Name[i])
						if err != nil {
							fmt.Println(config.CSI + config.Red + err.Error() + config.End)
							return
						}
					}
				} else if act.Env == "python" || act.Env == "Python" {
					for i := 0; i < len(act.Name); i++ {
						fmt.Println(config.CSI + config.Blue + ">Run " + act.Name[i] + config.End)
						err := config.Execute(act.Env, act.Name[i])
						if err != nil {
							fmt.Println(config.CSI + config.Red + err.Error() + config.End)
							return
						}

					}
				}
			} else {
				fmt.Println(config.CSI + config.Red + act.Action + " cant run" + config.End)
			}
		}
	}
}
func main() {
	kingpin.Parse()
	if !*doYaml && !*doAll && !*doHelp && !*doVersion {
		showSystem()
		showCpu()
		showTimeZone()
		showNet()
		showDisk()
		showMem()
		showNtp()
		showIptables()
		showDns()
		showIpmi()
	}
	if *doYaml {
		doYamlAction()
	}
	if *doAll {
		showSystem()
		showCpu()
		showTimeZone()
		showNet()
		showDisk()
		showMem()
		showNtp()
		showIptables()
		showDns()
		showIpmi()
		doYamlAction()
	}
	if *doHelp {
		usage()
	}
	if *doVersion {
		showVersion()
	}
}
