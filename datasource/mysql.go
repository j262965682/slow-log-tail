package datasource

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"slow-log-tail/config"
	"slow-log-tail/model"
	"time"
)

var orm *gorm.DB

func GetDB() *gorm.DB {
	if orm != nil {
		return orm
	}
	return InitDB()
}

func InitDB() *gorm.DB {
	mysqlConfig := config.GetConfig().Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", mysqlConfig.User, mysqlConfig.PassWord, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Dbname, mysqlConfig.Other)
	//fmt.Println(dsn)
	// MySQl 驱动程序提供了 一些高级配置 可以在初始化过程中使用
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		// 使用CreateBatchSize 选项初始化 GORM 时，所有的创建& 关联 INSERT 都将遵循该选项
		CreateBatchSize: 1000,
		// 全局模式：执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		PrepareStmt: true,
		// 注意 QueryFields 模式会根据当前 model 的所有字段名称进行 select。
		//QueryFields: true,
		Logger: logger.Default,
		// 表名加前缀和禁用复数表名
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀，`User`表为`t_users`
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	})

	if err != nil {
		log.Println("conn mysql fail", err)
		return nil
	}

	sqlDb, _ := db.DB()
	// 对于中小型 web 应用程序，我通常使用以下设置作为起点，然后根据负载测试结果和实际吞吐量级别进行优化。
	// SetMaxIdleConns: 设置空闲连接池中链接的最大数量
	sqlDb.SetMaxIdleConns(10)
	// SetMaxOpenConns: 设置打开数据库链接的最大数量
	sqlDb.SetMaxOpenConns(10)
	// SetConnMaxLifetime: 设置链接可复用的最大时间
	sqlDb.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&model.SlowLog{})

	return db

}
