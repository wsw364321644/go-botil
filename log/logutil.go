package log

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/wsw364321644/go-botil"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)
type myFormatter struct {
	logrus.TextFormatter
	terminalInitOnce sync.Once
	isTerminal bool
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// this whole mess of dealing with ansi color codes is required if you want the colored output otherwise you will lose colors in the log levels
	f.terminalInitOnce.Do(func() { f.init(entry) })
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	if f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows")){
		return []byte(fmt.Sprintf("[%s] - \x1b[%dm%s\x1b[0m - %s\n", entry.Time.Format(f.TimestampFormat), levelColor, strings.ToUpper(entry.Level.String()), entry.Message)), nil
	}else{
		return []byte(fmt.Sprintf("[%s] - %s - %s\n", entry.Time.Format(f.TimestampFormat), strings.ToUpper(entry.Level.String()), entry.Message)), nil
	}
}
func (f *myFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}

var FileLog = logrus.New()
var ConsoleLog = logrus.New()

var logList []*logrus.Logger=make([]*logrus.Logger,0)
var logsettings *LogSettings

type LogSettings struct{
	Name string
	Rotate bool
	Level string
}
func Init(settings *LogSettings ){
	logsettings=settings
	var f io.Writer
	var err error
	if logsettings.Rotate {
		f, err = rotatelogs.New(
			logsettings.Name+"-%Y%m%d"+".log",
			// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
			//rotatelogs.WithLinkName("llmmkr.log"),

			// WithRotationTime设置日志分割的时间
			rotatelogs.WithRotationTime(time.Hour*24),

			rotatelogs.WithClock(rotatelogs.UTC),
			// WithMaxAge和WithRotationCount二者只能设置一个，
			// WithMaxAge设置文件清理前的最长保存时间，
			// WithRotationCount设置文件清理前最多保存的个数。
			rotatelogs.WithMaxAge(time.Hour*24*90),
			//rotatelogs.WithRotationCount(maxRemainCnt),
		)
		if err != nil {
			log.Panicf("failed to create rotatelogs: %s", err)
		}
	}else {
		file, err := os.OpenFile(logsettings.Name+".log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Panicf("error opening file: %v", err)
		}
		_, err = file.Stat()
		if err != nil {
			log.Panicf("error opening file: %v", err)
		}
		f=file
	}

	FileLog.Out=f
	ConsoleLog.Out=os.Stdout

	if logsettings.Level==""{
		logsettings.Level="debug"
	}
	UpdateLevel(logsettings.Level)

	myformatter:=myFormatter{logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            false,
		DisableLevelTruncation: true,
	},sync.Once{},false}
	ConsoleLog.SetFormatter(&myformatter)
	FileLog.SetFormatter(&myformatter)
	logList=append(logList, FileLog,ConsoleLog)
}
func UpdateLevel(level string){
	lvl,err:=logrus.ParseLevel(level)
	if err==nil{
		FileLog.SetLevel(lvl)
		ConsoleLog.SetLevel(lvl)
	}
}

//func Println(v ...interface{}){
//	ConsoleLog.Println()
//	FileLog.Info(v)
//}
//func Printf(format string, args ...interface{}){
//	Log.Infof(format,args)
//}
func Panic(args ...interface{}){
	FileLog.Panic(args...)
}
func Panicf(format string, args ...interface{}){
	FileLog.Panicf(format,args...)
}
func Panicln(args ...interface{}){
	FileLog.Panicln(args...)
}


func Error(args ...interface{}){
	for _,log :=range logList{
		log.Error(ConvertArgs(args)...)
	}
}
func Errorf(format string, args ...interface{}){
	for _,log :=range logList{
		log.Errorf(format,ConvertArgs(args)...)
	}
}
func Errorln(args ...interface{}){
	for _,log :=range logList{
		log.Errorln(ConvertArgs(args)...)
	}
}

func ConvertArgs(args ...interface{})[]interface{}{
	newargs:=[]interface{}{}
	for _,arg:=range args{
		llerr,ok:=arg.(*botil.LLError)
		if ok{
			newargs=append(newargs, llerr.Message())
		}else{
			newargs=append(newargs,arg)
		}
	}
	return newargs
}

func Warn(args ...interface{}){
	for _,log :=range logList{
		log.Warn(args...)
	}
}
func Warnf(format string, args ...interface{}){
	for _,log :=range logList{
		log.Warnf(format,args...)
	}
}
func Warnln(args ...interface{}){
	for _,log :=range logList{
		log.Warnln(args...)
	}
}

func Info(args ...interface{}){
	for _,log :=range logList{
		log.Info(args...)
	}
}
func Infof(format string, args ...interface{}){
	for _,log :=range logList{
		log.Infof(format,args...)
	}
}
func Infoln(args ...interface{}){
	for _,log :=range logList{
		log.Infoln(args...)
	}
}

func Debug(args ...interface{}){
	for _,log :=range logList{
		log.Debug(args...)
	}
}
func Debugf(format string, args ...interface{}){
	for _,log :=range logList{
		log.Debugf(format,args...)
	}
}
func Debugln(args ...interface{}){
	for _,log :=range logList{
		log.Debugln(args...)
	}
}

func GetCronLog()*logrus.Logger{
	return FileLog
}