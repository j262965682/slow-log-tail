package parsefile

import (
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/percona/go-mysql/query"
	"github.com/pkg/errors"
	"log"
	"regexp"
	"slow-log-tail/config"
	"slow-log-tail/model"
	"time"
)

// ParseTail fileName "/var/log/nginx.log"  "D:\golangPro\slow-log-tail\slowlog\slow-log.log"
// func ParseTail(fileName string, db *gorm.DB, instance string, ignoreUser []string, loneQueryTime float64) {
func ParseTail(config *config.Config, sender model.Sender) {
	// tail.TailFile()函数开启goroutine去读取文件，通过channel格式的t.lines传递内容。
	t, err := tail.TailFile(config.SlowLog.Path, tail.Config{Follow: true, ReOpen: true, MustExist: false})
	if err != nil {
		err = errors.Wrap(err, "ERROR -> open file failed.") //如果文件不存在，会阻塞并打印Waiting for my.log to appear...，直到文件被创建
		log.Println(err)
	}
	var slowLog model.SlowLog
	var slowClear model.SlowLog
	var userAndHost model.UserAndHost
	var row model.Row
	var timestampValue int
	var dbValue string
	var signGroup bool //slow log group control mark
	var signSQL bool   //sql control mark
	var sql string
	var lineType string
	var lineValue string
	var fingerprint string

	timeReg := regexp.MustCompile(`:\s(.*)`)
	userReg := regexp.MustCompile(`\[(.*?)\]`)
	signSQL = false

	for line := range t.Lines {
		//fmt.Printf(".")
		//for {
		//	line, ok := <-t.Lines
		//	if !ok {
		//		fmt.Println("tail file close,fileName:", t.Filename)
		//		continue
		//	}
		//fmt.Println("signGroup:", signGroup)
		//fmt.Println("signSQL:", signSQL)
		//fmt.Println(line.Text)
		lineType = TypeOfLine(line.Text, signSQL)
		fmt.Printf(" %s.", lineType)
		switch lineType {
		case "TIME":
			//this line is the slowLog group of first,so clear struct and set sign true.
			signGroup = true //如果遇到 time 就强行清空结构体 和 sql连接标识
			signSQL = false

			slowLog = slowClear
			sql = ""

			slowLog.Time = timeRule(line.Text, timeReg)
			fmt.Printf(" %s ,%s-1", slowLog.Time, timeFormat()) //时间戳
		case "USER":
			if signGroup {
				userAndHost = userRule(line.Text, userReg)
				slowLog.User = userAndHost.User
				slowLog.Host = userAndHost.Host
				slowLog.ThreadId = userAndHost.ThreadId
				fmt.Printf(" ,%s-2", timeFormat()) //用户
			}
		case "ROW":
			if signGroup {
				row, err = rowRule(line.Text)
				if err = errors.Wrap(err, "ERROR -> ROW parse failed."); err != nil {
					log.Println(err)
				}
				slowLog.QueryTime = Round(row.QueryTime*1000, 3)
				slowLog.LockTime = Round(row.LockTime*1000, 3)
				slowLog.RowsSend = row.RowsSend
				slowLog.RowsExamined = row.RowsExamined
				fmt.Printf(" ,%s-3", timeFormat()) //执行时间
			}
		case "TIMESTAMP":
			if signGroup {
				timestampValue, err = timestampRule(line.Text)
				if err = errors.Wrap(err, "ERROR -> ROW parse failed."); err != nil {
					log.Println(err)
				}
				slowLog.Timestamp = timestampValue
				fmt.Printf(" ,%s-4", timeFormat()) //执行时间
			}
		case "DB":
			//if exist line of 'use xxx',then get it;if not exit ,then the value is the last.
			if signGroup {
				dbValue = dbRule(line.Text)
				//slowLog.Db = dbValue
				fmt.Printf(" ,%s-5", timeFormat()) //库名
			}
		case "SQL":
			if signGroup {
				signSQL = true
				lineValue = TrimString(line.Text)
				//fmt.Println("sql" + lineValue + ".")
				if len(lineValue) > 0 {
					sql = sql + " " + lineValue
					b := lineValue[len(lineValue)-1]
					if b == ';' {
						fingerprint = MaoHaoChange(query.Fingerprint(sql))
						sql = MaoHaoChange(sql)
						slowLog.Db = dbValue
						slowLog.Sql = sql
						slowLog.Env = config.SlowLog.Env
						slowLog.Instance = config.SlowLog.Instance
						slowLog.Fingerprint = fingerprint
						slowLog.Hash = Md532(fingerprint)
						fmt.Printf(" ,%s-6", timeFormat()) // sql
						signGroup = false
						signSQL = false
						//fmt.Println("  sql:", sql, ".")
						//fmt.Println(ignoreUser)
						// 判断慢查询阈值
						//fmt.Println("QueryTime is", slowLog.QueryTime)
						//fmt.Println("loneQueryTime is", loneQueryTime)
						if slowLog.QueryTime >= config.SlowLog.LongQueryTime {
							//fmt.Println("大于阈值")
							// 判断用户
							if !In(slowLog.User, config.SlowLog.IgnoreUser) {
								//send slow log to database
								//if err = slowLog.SendToDatabase(db); err != nil {
								//	fmt.Println(err)
								//}
								if err = sender.SendTo(&slowLog); err != nil {
									log.Println(err)
								}
								fmt.Printf(" ,%s-7", timeFormat()) // 发送
							}
						}
					}
				}
			}
		}
	}
}

func timeFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//
//func tailInitConfig() tail.Config {
//	config := tail.Config{
//		Location: &tail.SeekInfo{
//			Offset: 0,
//			Whence: 0,
//		},
//		Poll:      true,
//		ReOpen:    true,
//		MustExist: false,
//		Follow:    true,
//	}
//	return config
//}
//
//func handleTail(fileName string, config tail.Config) (*tail.Tail, error) {
//	tails, err := tail.TailFile(fileName, config)
//	if err != nil {
//		return nil, errors.Wrap(err, "ERROR-> get file tail handle failed")
//	}
//	return tails, nil
//}
//
//func GetFileTailHandle(fileName string) (tails *tail.Tail, err error) {
//	config := tailInitConfig()
//
//	tails, err = handleTail(fileName, config)
//	if err != nil {
//		return nil, err
//	}
//	return tails, nil
//}
//
//func RuleParsinge(tails *tail.Tail) {
//	var line *tail.Line
//	var ok bool
//	for {
//		// 通过道获取到每条行管数据
//		line, ok = <-tails.Lines
//		fmt.Println("走这里了")
//		if !ok {
//			fmt.Println("tail file close,fileName:", tails.Filename)
//			time.Sleep(1 * time.Second)
//			continue
//		}
//		fmt.Println("line:", line)
//	}
//}
