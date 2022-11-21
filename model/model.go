package model

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"time"
)

//type LineType int
//
//const (
//	TIME LineType = iota
//	DB
//	USER
//	ROW
//	TIMESTAMP
//	SQL
//)

type SlowLog struct {
	ID           int64   `gorm:"column:id;primaryKey;autoIncrement:true"`
	Instance     string  `gorm:"column:instance;type:varchar(30);index:idx_db_instance,priority:2"`
	Db           string  `gorm:"column:db;type:varchar(30);index:idx_db_instance,priority:1"`
	Time         string  `gorm:"column:time;type:varchar(35);index:idx_time"`
	User         string  `gorm:"column:user;type:varchar(30);index:idx_user"`
	Host         string  `gorm:"column:host;type:varchar(30);"`
	ThreadId     string  `gorm:"column:thread_id;type:varchar(15);"`
	QueryTime    float64 `gorm:"column:query_time;type:float"`
	LockTime     float64 `gorm:"column:lock_time;type:float"`
	Timestamp    int     `gorm:"column:timestamp;type:int"`
	RowsSend     int     `gorm:"column:rows_send;type:int"`
	RowsExamined int     `gorm:"column:rows_examined;type:int"`
	Sql          string  `gorm:"column:sql;type:text"`
	//Fingerprint  string  `gorm:"column:finger_print;type:text;index:idx_finger_print"`
	Hash      string `gorm:"column:hash;type:varchar(32);index:idx_hash"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserAndHost struct {
	User     string
	Host     string
	ThreadId string
}

type Row struct {
	QueryTime    float64
	LockTime     float64
	RowsSend     int
	RowsExamined int
}

type SlowLogRepository struct {
	db *gorm.DB
}

//func (s *SlowLog) SendToDatabase(db *gorm.DB) (err error) {
//	if err = db.Create(s).Error; err != nil {
//		errors.Wrap(err, "ERROR -> insert DB failed.")
//		return err
//	}
//	return nil
//}

type Mysql struct {
	Db *gorm.DB
}

func (m *Mysql) SendTo(s *SlowLog) error {
	var err error
	if err = m.Db.Create(s).Error; err != nil {
		errors.Wrap(err, "ERROR -> insert DB failed.")
		return err
	}
	return nil
}

type Post struct {
	Host string `yaml:"host"`
	Url  string `yaml:"url"`
}

func (m *Post) SendTo(s *SlowLog) error {
	req := fasthttp.AcquireRequest()   //获取Request连接池中的连接
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")
	req.SetRequestURI(m.Url)
	req.SetHost(m.Host)
	byteJson, _ := json.Marshal(s)
	req.SetBody(byteJson)
	resp := fasthttp.AcquireResponse()             //获取Response连接池中的连接
	defer fasthttp.ReleaseResponse(resp)           // 用完需要释放资源
	if err := fasthttp.Do(req, resp); err != nil { //发送请求
		errors.Wrap(err, "HTTP send error")
		return err
	}
	return nil
	//b := resp.Body()
	//// resp.Body()
	//log.Info(b)
}

type Sender interface {
	SendTo(s *SlowLog) error
}
