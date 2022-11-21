package main

import (
	"log"
	"slow-log-tail/config"
	"slow-log-tail/datasource"
	"slow-log-tail/model"
	"slow-log-tail/parsefile"
)

func main() {
	log.Println("Start parse config file and Connect to the specified database.")
	config := config.GetConfig()
	var sender model.Sender
	switch config.OutPutType {
	case "mysql":
		//fmt.Println("mysql")
		//解析配置文件 、 连接数据库
		datasource.InitDB()
		sender = &model.Mysql{Db: datasource.GetDB()}
	case "post":
		//fmt.Println("post")
		sender = &model.Post{
			Host: config.Post.Host,
			Url:  config.Post.Url,
		}
	}
	log.Println("Start parse the slow log file.")
	//fmt.Println(config.GetConfig().SlowLog.IgnoreUser)
	parsefile.ParseTail(config, sender)
}
