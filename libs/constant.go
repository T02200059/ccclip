package libs

import (
	"runtime"
	"time"
)

const CloudURL = "http://192.168.1.139:22122"
const CloudPort = ":22122"

const CopyCollectorInterval = 3 * time.Second
const PingPongInterval = 10 * time.Second
const NetworkErrInterval = 10 * time.Second

const SuffixAsync = "/async"

type UserCode int

const OriginUserYTB UserCode = 100

const (
	SysUnknown UserCode = 0
	SysLinux   UserCode = 1
	SysWindows UserCode = 2
	SysDarwin  UserCode = 3
	SysAndroid UserCode = 4
)

func GetPlatformCode() (p UserCode) {
	println("Your os is: ", runtime.GOOS)

	switch runtime.GOOS {
	case "linux":
		p = SysLinux
	case "windows":
		p = SysWindows
	case "darwin":
		p = SysDarwin
	default:
		p = SysUnknown
	}
	return
}

func DecodeUser(n UserCode) (p UserCode) {
	c := n / 100 % 10
	p = UserCode(c * 100)
	return
}

func DecodePlatform(n UserCode) (p UserCode) {
	c := n % 10
	p = UserCode(c)
	return
}
