package parsefile

import (
	"context"
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

func WatchLogFile(config *config.Config, sender model.Sender, ctx context.Context, keyChan chan<- string) {
	fmt.Println("begin goroutine watch log file ", config.SlowLog.Path)
	tailFile, err := tail.TailFile(config.SlowLog.Path, tail.Config{
		//文件被移除或被打包，需要重新打开
		ReOpen: true,
		//实时跟踪
		Follow: true,
		//如果程序出现异常，保存上次读取的位置，避免重新读取
		//Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		//支持文件不存在
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	handleLines(config, sender, tailFile, ctx)

	defer func() {
		if errCover := recover(); errCover != nil {
			fmt.Println("goroutine watch ", config.SlowLog.Path, " panic")
			fmt.Println(errCover)
			keyChan <- config.SlowLog.Path
		}
	}()
}

func handleLines(config *config.Config, sender model.Sender, t *tail.Tail, ctx context.Context) {
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
	var err error

	timeReg := regexp.MustCompile(`:\s(.*)`)
	userReg := regexp.MustCompile(`\[(.*?)\]`)
	signSQL = false

	for {
		select {
		case line, ok := <-t.Lines:
			if !ok {
				fmt.Println("tail file close,fileName:", t.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			//fmt.Println("signGroup:", signGroup)
			//fmt.Println("signSQL:", signSQL)
			//fmt.Println(line.Text)
			lineType = TypeOfLine(line.Text, signSQL)
			//fmt.Println("lineType:", lineType)
			switch lineType {
			case "TIME":
				//this line is the slowLog group of first,so clear struct and set sign true.
				signGroup = true //如果遇到 time 就强行清空结构体 和 sql连接标识
				signSQL = false

				slowLog = slowClear
				sql = ""

				slowLog.Time = timeRule(line.Text, timeReg)
				fmt.Printf("\n %s ,%s-1", slowLog.Time, timeFormat()) //时间戳
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
		case <-ctx.Done():
			fmt.Println("receive main gouroutine exit msg")
			fmt.Println("watch log file ", config.SlowLog.Path, " goroutine exited")
			return
		}
	}
}
