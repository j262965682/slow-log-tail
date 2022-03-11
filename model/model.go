package model

import (
	"github.com/pkg/errors"
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

func (s *SlowLog) SendToDatabase(db *gorm.DB) (err error) {
	if err = db.Create(s).Error; err != nil {
		errors.Wrap(err, "ERROR -> insert DB failed.")
		return err
	}
	return nil
}
