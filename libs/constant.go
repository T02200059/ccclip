package libs

import (
	"runtime"
	"time"
)

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
		p = 1
	case "windows":
		p = 2
	case "darwin":
		p = 3
	default:
		p = 0
	}
	return
}
