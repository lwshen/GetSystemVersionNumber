// golang 获取系统版本号(windows, linux, macOS)
// golang get system version number(windows, linux, macOS)
// thx https://gist.github.com/flxxyz/ae3ef071dc4ffb0c55daedc7f0740611

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type OS struct {
	Name    string
	Version string
	Arch    string
}

func NewOS() *OS {
	return &OS{
		Name:    runtime.GOOS,
		Version: "unknown",
		Arch:    runtime.GOARCH,
	}
}

// GetWindowsVersion 获取windows版本号
func GetWindowsVersion() (version string) {
	versionNumbers := map[string]string{
		`5\.0`:  "2000",
		`5\.1`:  "XP",
		`5\.2`:  "Server 2003",
		`6\.0`:  "Server 2008",
		`6\.1`:  "Server 2008 R2",
		`6\.2`:  "Server 2012",
		`6\.3`:  "Server 2012 R2",
		`10\.0`: "10",
	}

	//win10VersionNumbers := map[string]string{
	//	`10\.0\.14300`: "Server 2016",
	//	`10\.0\.14393`: "Server 2016",
	//	`10\.0\.16299`: "Server 2016",
	//	`10\.0\.17134`: "Server 2016",
	//	`10\.0\.17677`: "Server 2019",
	//	`10\.0\.17763`: "Server 2019",
	//	`10\.0\.18362`: "Server 2019",
	//	`10\.0\.18363`: "Server 2019",
	//}

	cmd := exec.Command("cmd.exe")
	out, _ := cmd.StdoutPipe()
	buffer := bytes.NewBuffer(make([]byte, 0))
	cmd.Start()
	buffer.ReadFrom(out)
	str, _ := buffer.ReadString(']')
	cmd.Wait()
	for key, _ := range versionNumbers {
		re := regexp.MustCompile(`Microsoft Windows \[[\s\S]* ` + key + `\.([0-9]+).?[0-9]*\]`)
		if re.MatchString(str) {
			if versionNumbers[key] != "10" {
				version = versionNumbers[key]
			} else {
				versionNumber := re.FindStringSubmatch(str)
				if len(versionNumber) > 1 {
					if Str2Int(versionNumber[1]) <= 17134 {
						version = "Server 2016"
					} else {
						version = "Server 2019"
					}
				}
			}

			return
		}
	}

	return
}

// GetLinuxVersion 获取linux版本号
func GetLinuxVersion() (name, version string) {
	if ok, _ := PathExists("/etc/os-release"); ok {
		cmd := exec.Command("cat", "/etc/os-release")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		content, err := ioutil.ReadAll(stdout)
		if err == nil {
			id := regexp.MustCompile(`\nID="?(.*?)"?\n`).FindStringSubmatch(string(content))
			if len(id) > 1 {
				name = id[1]
			}

			versionId := regexp.MustCompile(`VERSION_ID="?([.0-9]+)"?\n`).FindStringSubmatch(string(content))
			if len(versionId) > 1 {
				version = versionId[1]
			}
		}
	}

	return
}

// GetMacVersion 获取mac版本号
func GetMacVersion() (version string) {
	cmd := exec.Command("sw_vers", "-productVersion")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		return
	}
	version = string(content)
	version = strings.TrimSpace(version)
	return
}

// PathExists 检查路径存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func Str2Int(value string) (val int) {
	val, _ = strconv.Atoi(value)

	return
}

func main() {
	os := NewOS()

	switch runtime.GOOS {
	case "windows":
		os.Version = GetWindowsVersion()
	case "linux":
		os.Name, os.Version = GetLinuxVersion()
	case "darwin":
		os.Version = GetMacVersion()
	}

	log.Println(os)
}
