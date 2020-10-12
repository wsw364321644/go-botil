package botil

import (
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"syscall"
)

func RunMeElevated(args ...string) error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	conargs := strings.Join(args, " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(conargs)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		return err
	}
	return nil
}

func AmAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}