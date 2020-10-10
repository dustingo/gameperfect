package config

// 硬件信息包，包括cpu，system，timezone，memory，net,ntpd,iptables,ipmitool

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
)

// CpuInfo cpu基础信息
type CpuInfo struct {
	ModelName  string
	CpuMHz     string
	PhysicalID []string
	CpuCores   string
	Processor  []string
}

// SystemInfo 系统版本，内核信息，systemd版本，时区
type SystemInfo struct {
	OSName    string
	OSRelease string
	OSKernel  string
	OSArch    string
}

// Timezone  "timedatectl结构体
type Timezone struct {
	LocalTime string
	Zone      string
}

// NetFaceInfo 网卡信息
type NetFaceInfo struct {
	//NetInfo []NetInfo
	Name      []string
	Speed     []string
	Ipaddress []string
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	MemTotal     string
	MemAvailable string
	MemSwap      string
}

// SystemInfo 基本系统信息
func (system *SystemInfo) GetSystem() {
	fversion, _ := os.Open("/proc/version")
	versionData, _ := ioutil.ReadAll(fversion)
	verData := strings.Split(string(versionData), " ")
	system.OSKernel = verData[2]
	system.OSArch = strings.Split(verData[2], ".")[4]
	redHat, err := os.Open("/etc/redhat-release")
	if err == nil {
		releaseData, _ := ioutil.ReadAll(redHat)
		redData := strings.Split(string(releaseData), " ")
		system.OSName = redData[0]
		system.OSRelease = redData[3]
	} else { // ubuntu
		fissue, _ := os.Open("/etc/issue")
		defer fissue.Close()
		issueData, _ := ioutil.ReadAll(redHat)
		issData := strings.Split(string(issueData), " ")
		system.OSName = issData[0]
		system.OSRelease = issData[1]
	}
	fversion.Close()
	redHat.Close()

}

// GetCpu 获取cpu信息
func (cpu *CpuInfo) GetCpu() {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	//data, err := ioutil.ReadAll(f)
	//fmt.Println(string(data))
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		//parts := strings.Fields(line)
		if strings.HasPrefix(line, "processor") {
			cpu.Processor = append(cpu.Processor, strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.HasPrefix(line, "model name") {
			cpu.ModelName = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "cpu MHz") {
			cpu.CpuMHz = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "physical id") {
			cpu.PhysicalID = append(cpu.PhysicalID, strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.HasPrefix(line, "cpu cores") {
			cpu.CpuCores = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
}

// GetTimeZone 获取时区信息
func (tzone *Timezone) GetTimeZone() {
	cmd := exec.Command("timedatectl")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range strings.Split(strings.Replace(out.String(), " ", "", -1), "\n") {
		if len(item) == 0 {
			continue
		}
		if strings.HasPrefix(item, "Localtime") {
			tzone.LocalTime = strings.SplitN(item, ":", 2)[1]
		}
		if strings.HasPrefix(item, "Timezone") {
			tzone.Zone = strings.SplitN(item, ":", 2)[1]
		}
	}

}

// GetNetInfo 获取网卡信息
func (netface *NetFaceInfo) GetNetInfo() {
	face, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(face); i++ {
		var out bytes.Buffer
		if face[i].Name == "lo" || !strings.HasPrefix(face[i].Flags.String(), "up") {
			continue
		}
		cmd := exec.Command("ethtool", face[i].Name)
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		netface.Name = append(netface.Name, face[i].Name) // ok
		sliceOut := strings.Split(strings.Replace(out.String(), " ", "", -1), "\n")
		for j := 0; j < len(sliceOut); j++ {
			if strings.HasPrefix(strings.TrimSpace(sliceOut[j]), "Speed") { // 过滤seppd:xxxx
				//fmt.Println(face[i].Name)
				byName, err := net.InterfaceByName(face[i].Name)
				if err != nil {
					fmt.Println(err)
					return
				}
				address, err := byName.Addrs()
				netface.Speed = append(netface.Speed, strings.Split(strings.TrimSpace(sliceOut[j]), ":")[1])
				if len(address) == 0 {
					netface.Ipaddress = append(netface.Ipaddress, "nil.nil.nil.nil")
				} else {
					netface.Ipaddress = append(netface.Ipaddress, address[0].String())
				}

			}
		}

	}
}

// GetMemInfo 获取内存信息
func (mem *MemoryInfo) GetMemInfo() {
	file, err := os.Open(ProcFilepath("meminfo"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "MemTotal") {
			mem.MemTotal = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "MemAvailable") {
			mem.MemAvailable = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "SwapCached") {
			mem.MemSwap = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

}

// GetDnsInfo 获取dns信息
func GetDnsInfo() {
	file, err := ioutil.ReadFile(EtcFilepath("resolv.conf"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(file))
}

// GetNtpInfo 获取NTP信息
func GetNtpInfo() {
	err := Execute("sh", "-c", "./tools/ntp.sh")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// GetIptables 获取防火墙信息
func GetIptables() {
	err := Execute("sh", "-c", "./tools/iptables.sh")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// GetIpmiInfo 获取ipmi信息
func GetIpmiInfo() {
	if YumCheck("ipmitool") {
		err := Execute("/usr/bin/ipmitool", "sel", "list")
		if err != nil {

			return
		}
	} else {
		err := Execute("yum", "-y", "install", "ipmitool")
		if err != nil {
			fmt.Println(CSI + Red + err.Error() + End)
			return
		}
		err = Execute("/usr/bin/ipmitool", "sel", "list")
		if err != nil {

			return
		}
	}
}
