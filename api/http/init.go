package http

import (
	"fmt"
	"os/exec"
)

func init() {
	pypath, err := lookupPythonPath()
	if err != nil {
		fmt.Printf("cannot find python error:%s\n", err)
		return
	}
	fmt.Printf("use python(%s)\n", pypath)
	err = installUiautomator()

	if err != nil {
		fmt.Printf("run `pip3 install uiautomator2` error:%s", err)
		return
	}
}

func lookupPythonPath() (pypath string, err error) {
	pypath, err = exec.LookPath("python3")
	if err == nil {
		return
	}
	pypath, err = exec.LookPath("python")
	if err == nil {
		return
	}
	return
}
func installUiautomator() error {
	return exec.Command("pip3", "install", "uiautomator2").Run()
}
