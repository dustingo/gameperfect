package config

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
)

type DomainUrl struct {
	Domain []string `yaml:"resolve"`
}

var wg sync.WaitGroup

// 执行get请求
func getHttps(url string) { // url,result,StatusCode
	defer wg.Done()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	statusCode := response.StatusCode
	if statusCode == 200 {
		fmt.Printf("\033[1;36;40m%s\033[0m %s \033[1;36;40m%s\033[0m %s \033[1;36;40m%s\033[0m %d\n",
			"URL:", url, "Result:", "Success", "StatusCode:", response.StatusCode)
	} else {
		fmt.Printf("\033[1;36;40m%s\033[0m %s \033[1;36;40m%s\033[0m %s \033[1;36;40m%s\033[0m %d\n",
			"URL:", url, "Result:", "Failed", "StatusCode:", response.StatusCode)
	}

}

// ResolveDns 解析dns.yaml，执行域名解析和get
func ResolveDns() {
	yamlConfig := new(DomainUrl)
	yamlByte, err := ioutil.ReadFile("./config/dns.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(yamlByte) == 0 {
		fmt.Printf(CSI+Red+"%s"+End+"\n", "dns.yaml is null")
		return
	}
	err = yaml.Unmarshal(yamlByte, yamlConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(yamlConfig.Domain); i++ {
		if strings.HasPrefix(yamlConfig.Domain[i], "https") {
			wg.Add(1)
			go getHttps(yamlConfig.Domain[i])
		} else {
			ipaddrs, err := net.LookupIP(yamlConfig.Domain[i])
			if err != nil {
				fmt.Printf(CSI+Red+"%v"+End+"\n", err)
				continue
			}
			fmt.Printf(CSI+Blue+"%s: "+End+"%s"+"\n", "Domain", yamlConfig.Domain[i])
			for i := 0; i < len(ipaddrs); i++ {
				fmt.Printf(CSI+Blue+"IP%d: "+End+"%s"+"\n", i, ipaddrs[i])
			}
		}
	}
	wg.Wait()
}
