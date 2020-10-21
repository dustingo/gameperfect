package config

// 通用包

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var (
	ProcPath          = "/proc"
	IgnoreMountPoints = "^/(boot|dev|proc|sys|run|var/lib/docker/.+)($|/)"
	EtcPath           = "/etc"
)

// color
const (
	CSI   = "\033["
	Reset = CSI + "m"
	Red   = CSI + "31;40m"
	Blue  = CSI + "36;40m"
	Green = CSI + "32;40m"
	End   = CSI + "0m"
)

// ProcFilepath /proc 组合路径
func ProcFilepath(name string) string {
	return filepath.Join(ProcPath, name)
}

// EtcFilepath /etc 组合路径
func EtcFilepath(name string) string {
	return filepath.Join(EtcPath, name)
}

var sg sync.WaitGroup

func asyncPrint(reader io.ReadCloser) {
	defer sg.Done()
	scan := bufio.NewScanner(reader)
	for scan.Scan() {
		line := scan.Text()
		fmt.Println(line)
	}
}

// Execute 调用shell
func Execute(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s......\n", err.Error())
		return err
	}
	sg.Add(2)
	go asyncPrint(stdout)
	go asyncPrint(stderr)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for command execution: %s......\n", err.Error())
		return err
	}
	sg.Wait()
	time.Sleep(time.Second * 2)
	return nil
}

// YumCheck 检查是否已安装
func YumCheck(name string) bool {
	env := os.Environ()
	procAttr := &os.ProcAttr{
		Env: env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}
	fmt.Println(CSI + Blue + name + "版本：" + End)
	newProcess, err := os.StartProcess("/usr/bin/rpm", []string{"rpm", "-q", name}, procAttr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	processState, err := newProcess.Wait()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return processState.Success()

}

// 判断目录是否存在
func PathExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(CSI + Blue + ">>create dir: " + name + End)
			return false
		}
		return true
	}
	return true
}

// Mkdir 创建目录
func Mkdir(name string, perm os.FileMode) {
	err := os.MkdirAll(name, perm)
	if err != nil {
		fmt.Println(CSI + Red + err.Error() + End)
	}
}
