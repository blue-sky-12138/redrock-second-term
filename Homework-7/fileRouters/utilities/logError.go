package utilities

import (
	"fmt"
	"log"
	"os"
)

var LogFile *log.Logger

//初始化日志输出
func LogInit() {
	File, err := os.OpenFile("./logFile.md", os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		panic("打开日志失败")
	}
	logFile := log.New(File, "[GIN-debug]", log.Ldate|log.Ltime)
	LogFile = logFile
}

//输出错误日志
func LogError(detail string, err error) {
	detail = "[" + detail + "]"
	log.Printf(detail, err, "\n")

	LogFile.Println(detail, err)
}
