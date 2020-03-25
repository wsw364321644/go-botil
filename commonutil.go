package botil

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"github.com/satori/go.uuid"
	"os"
	"reflect"
	"strings"
)


func GetArrayElemWithDefault(array interface{},index int,defaultValue interface{}) interface{}{
	v := reflect.ValueOf(array)
	if v.Kind() == reflect.Slice &&v.IsValid()&&v.Len()>0{
		return v.Index(0).Interface()
	} else{
		return defaultValue
	}
}

const (
	UID = 501
	GUID = 100
)
//func StartProcess (file string) {
//	var cred =  &syscall.Credential{ UID, GUID, []uint32{} }
//	// the Noctty flag is used to detach the process from parent tty
//	var sysproc = &syscall.SysProcAttr{  Credential:cred, Noctty:true }
//	var attr = os.ProcAttr{
//		Dir: ".",
//		Env: os.Environ(),
//		Files: []*os.File{
//			os.Stdin,
//			nil,
//			nil,
//		},
//		Sys:sysproc,
//
//	}
//	process, err := os.StartProcess(file, []string{"/bin/sleep", "100"}, &attr)
//	if err == nil {
//
//		// It is not clear from docs, but Realease actually detaches the process
//		err = process.Release();
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//
//	} else {
//		fmt.Println(err.Error())
//	}
//}

func CloseProcess(processname string){
	processes,err:=ps.Processes()
	if(err!=nil){
		fmt.Println(err)
		return
	}
	for _,process :=range processes{
		if(!strings.HasPrefix(process.Executable(),processname)){
			continue
		}
		p,err:=os.FindProcess(process.Pid())
		if(err!=nil){
			fmt.Println(err)
		}else{
			fmt.Printf("close pid:%d \n",process.Pid())
			p.Kill()
		}
	}
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func GetRandomID() string {
	id:=uuid.NewV4()
	return id.String()
}