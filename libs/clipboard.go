package libs

import (
	"github.com/atotto/clipboard"
	"runtime"
)

func ReadClipboard() (result string, err error) {
	switch runtime.GOOS {
	case "linux":
		result, err = clipboard.ReadAll()
		if err != nil {
			return
		}
	case "windows":

	case "darwin":

	}
	return
}

func WriteClipboard(payload string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = clipboard.WriteAll(payload)
		if err != nil {
			return
		}
	case "windows":

	case "darwin":

	}
	return
}
