package botil

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var readHellper *ReadHelper
func init(){
	readHellper=NewReadHelper()
	go func(){
		for {
			str, err := Scanfln()
			if(err!=nil){
				fmt.Println("read stdin err encountered:",err)
				return
			}
			readHellper.rwmutex.Lock()
			if readHellper.preadchan!=nil{
				*readHellper.preadchan<-str
			}
			readHellper.rwmutex.Unlock()
		}
	}()
}

type ReadHelper struct{
	preadchan *chan string
	rwmutex sync.RWMutex
}
func NewReadHelper()*ReadHelper {
	return new(ReadHelper)
}
func prepareRead(){
	readchan:=make(chan string,10)
	readHellper.preadchan=&readchan
}

func endRead(){
	readHellper.rwmutex.Lock()
	close(*readHellper.preadchan)
	readHellper.preadchan=nil
	readHellper.rwmutex.Unlock()
}
func Scanfln()(string,error){
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		return line,nil
	}
	if err := scanner.Err(); err != nil {
		return "",err
	}
	return "",errors.New("no input")
}

type CheckFunc func(string) bool
func CheckedScanfln(hint string,checkFunc CheckFunc)string{
	prepareRead()
	for{
		fmt.Print(hint)
		str:=<-*readHellper.preadchan
		if(checkFunc(str)){
			endRead()
			return str
		}
	}
}

func GetScanBoolFlag(printstring string,defaultbool bool) bool{
	str,flag := " ",defaultbool
	prepareRead()
	for {
		fmt.Print(printstring)
		str=<-*readHellper.preadchan
		if (strings.EqualFold(str, "y") || strings.EqualFold(str, "yes")) {
			flag=true
		}else if(strings.EqualFold(str, "n") || strings.EqualFold(str, "no")){
			flag=false
		}else if(str == ""){
		}else{
			continue
		}
		break;
	}
	endRead()
	return flag
}

func ReadLine() string{
	prepareRead()
	str:=<-*readHellper.preadchan
	endRead()
	return str
}