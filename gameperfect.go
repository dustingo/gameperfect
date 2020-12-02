package main

import (
	"bytes"
	"fmt"
	"gameperfect/config"
	"gopkg.in/alecthomas/kingpin.v2"
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
	choseMod  = kingpin.Flag("mode", "chose witch mode you want to run,must use with --yaml").String()
)
var (
	version = "1.0.1"
)

func usage() {
	fmt.Println(">>目录：config 存放yaml配置文件和第三方脚本配置")
	fmt.Println(">>目录:tools 存放自定义脚本")
	fmt.Println(">>config.yaml: ")
	fmt.Println("  >1.modeYum: action仅为'install'时，执行安装")
	fmt.Println("  >2.modeDir: action为'chown'时，根据para [用户.属组],更改目录属性，不存在则创建; action为'chmod'时，根据perm[]执行更改目录|文件权限")
	fmt.Println("  >3.modeScripts: action为'run'且匹配'hostname' 时，根据env [shell|python]和name [path + scripts_name]执行自定义脚本")
	fmt.Println(">4:modeService status 为active和inactive，active下面的服务都会被启动，且设置为开机启动，inactive下面的服务都设置为关闭，且开机不启动")
	fmt.Println(config.CSI + config.Red + "ATTENTION！注意：" + config.End)
	fmt.Println(config.CSI + config.Red + "the scripts should add 'x' permission" + config.End)
	fmt.Println(config.CSI + config.Red + "如果scripts是在win编辑的，在linux下记得使用dos转换" + config.End)
	fmt.Println(config.CSI + config.Red + "perm注意格式，如：0755 | 0600" + config.End)
	fmt.Println(config.CSI + config.Red + "host匹配服务器hostname时，只在匹配的服务器执行，当host为'all' 或 host为 '' 空时，任何服务器都执行" + config.End)
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
	if *choseMod == "" {
		config.DoYumAction(yamlConfig)
		config.DoDirAction(yamlConfig)
		config.DoScriptAction(yamlConfig)
		config.DoServiceAction(yamlConfig)
	} else if *choseMod == "yum" {
		config.DoYumAction(yamlConfig)
	} else if *choseMod == "dir" {
		config.DoDirAction(yamlConfig)
	} else if *choseMod == "scripts" {
		config.DoScriptAction(yamlConfig)
	} else if *choseMod == "services" {
		config.DoServiceAction(yamlConfig)
	} else {
		fmt.Println(config.CSI + config.Red + "parameter err" + config.End)
	}
}
func main() {
	kingpin.Parse()
	if !*doYaml && !*doAll && !*doHelp && !*doVersion && *choseMod == "" {
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
