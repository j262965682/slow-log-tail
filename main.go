package main

import (
	"log"
	"slow-log-tail/config"
	"slow-log-tail/datasource"
	"slow-log-tail/parsefile"
)

func main() {
	log.Println("Start parse config file and Connect to the specified database.")
	//解析配置文件 、 连接数据库
	datasource.InitDB()
	log.Println("Start parse the slow log file.")
	//fmt.Println(config.GetConfig().SlowLog.IgnoreUser)
	parsefile.ParseTail(config.GetConfig().SlowLog.Path, datasource.GetDB(), config.GetConfig().SlowLog.Instance, config.GetConfig().SlowLog.IgnoreUser)
}
