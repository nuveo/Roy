package logSys

import (
	"fmt"
	"path/filepath"
	"runtime"
)

var debugMode = false
var colors = true

func SetDebugMode(v bool) {
	debugMode = v
}

func logMsg(msg ...interface{}) (ret string) {
	_, fn, line, _ := runtime.Caller(2)
	fn = filepath.Base(fn)

	ret = fmt.Sprintf("[error] %s:%d ", fn, line)
	ret += fmt.Sprint(msg...)

	return
}

func Println(msg ...interface{}) {
	fmt.Println(logMsg(msg...))
}

func Print(msg ...interface{}) {
	fmt.Print(logMsg(msg...))
}
