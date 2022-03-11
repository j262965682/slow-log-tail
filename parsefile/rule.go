package parsefile

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"slow-log-tail/model"
	"strconv"
	"strings"
)

func TypeOfLine(line string, sqlControl bool) (lineType string) {
	if sqlControl {
		return "SQL"
	} else {
		//fmt.Println(line)
		//fmt.Printf("line length is %v \n", len(line))
		prefixStr := line[:3]
		switch prefixStr {
		case "# T":
			lineType = "TIME"
		case "# U":
			lineType = "USER"
		case "# Q":
			lineType = "ROW"
		case "SET":
			lineType = "TIMESTAMP"
		case "use":
			lineType = "DB"
		default:
			lineType = "SQL"
		}
		return lineType
	}
}

func timeRule(line string, reg *regexp.Regexp) (time string) {
	// |# Time: 2022-02-23T02:00:02.478019+08:00
	line = TrimString(line)
	many := reg.FindAllStringSubmatch(line, -1)
	return many[0][1]
}

func userRule(line string, reg *regexp.Regexp) (user model.UserAndHost) {
	// |# User@Host: root[root] @ localhost []  Id: 1895435
	line = TrimString(line)
	many := reg.FindAllStringSubmatch(line, -1)
	user.User = many[0][1]
	user.Host = many[1][1]
	user.ThreadId = line[strings.LastIndex(line, ":")+2:]
	return user
}

func rowRule(line string) (row model.Row, err error) {
	// |# Query_time: 0.016805  Lock_time: 0.000184 Rows_sent: 1  Rows_examined: 0
	line = TrimString(line)
	many := strings.Split(line, " ")
	if row.QueryTime, err = strconv.ParseFloat(many[2], 64); err != nil {
		errors.Wrap(err, "ERROR -> parse row QueryTime failed")
		return model.Row{}, err
	}
	if row.LockTime, err = strconv.ParseFloat(many[5], 64); err != nil {
		errors.Wrap(err, "ERROR -> parse row LockTime failed")
		return model.Row{}, err
	}
	if row.RowsSend, err = strconv.Atoi(many[7]); err != nil {
		errors.Wrap(err, "ERROR -> parse row RowsSend failed")
		return model.Row{}, err
	}
	if row.RowsExamined, err = strconv.Atoi(many[10]); err != nil {
		errors.Wrap(err, "ERROR -> parse row RowsExamined failed")
		return model.Row{}, err
	}
	return row, nil
}

func timestampRule(line string) (timestamp int, err error) {
	// |SET timestamp=1645552804;
	line = TrimString(line)
	many := strings.Split(line, "=")
	timestampStr := many[1]
	timestampStr = timestampStr[0 : len(timestampStr)-1]
	return strconv.Atoi(timestampStr)
}

func dbRule(line string) string {
	// |use db_item;
	line = TrimString(line)
	many := strings.Split(line, " ")
	dbStr := many[1]
	dbStr = dbStr[0 : len(dbStr)-1]
	return dbStr

}

func TrimString(str1 string) string {
	str1 = strings.Replace(str1, "\r", "", -1)
	str1 = strings.Replace(str1, "\n", "", -1)
	return strings.TrimSpace(str1)
}

func parseLine(line string) *model.SlowLog {
	// TODO
	return nil
}

func Md532(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func In(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}
