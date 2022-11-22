package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

var (
	_config *Config
)

type Config struct {
	OutPutType string   `yaml:"outPutType"`
	Mysql      *Mysql   `yaml:"mysql"`
	Post       *Post    `yaml:"post"`
	SlowLog    *SlowLog `yaml:"slowLog"`
}

type Mysql struct {
	Host      string `yaml:"host"`
	Port      uint32 `yaml:"port"`
	User      string `yaml:"user"`
	PassWord  string `yaml:"passWord"`
	Dbname    string `yaml:"dbName"`
	TableName string `yaml:"tableName"`
	Other     string `yaml:"other"`
}

type SlowLog struct {
	Path          string   `yaml:"path"`
	Env           string   `yaml:"env"`
	Instance      string   `yaml:"instance"`
	IgnoreUser    []string `yaml:"ignoreuser"`
	LongQueryTime float64  `yaml:"longquerytime"`
}

type Post struct {
	Host string `yaml:"host"`
	Url  string `yaml:"url"`
}

//初始化全局配置
func init() {
	viper.SetConfigName("slow_config") //指定配置文件的文件名称(不需要制定配置文件的扩展名)
	//设置配置文件类型
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config") // 设置配置文件的搜索目录
	viper.AddConfigPath("../../config")
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "read config file wrong."))
	}

	err := viper.Unmarshal(&_config) // 将配置信息绑定到结构体上
	if err != nil {
		panic(errors.Wrap(err, "Unmarshalling config file wrong."))
	}

	if err := checkMysqlConfig(_config); err != nil {
		panic(errors.Wrap(err, "mysqlConfig message wrong."))
	}

	//fmt.Println(_config)
	//viper.WatchConfig()
	////可以通过https://fsnotify.org 监听config文件变化更新配置信息
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("配置发生变更：", e.Name)
	//})
}

//获取全局配置
func GetConfig() *Config {
	return _config
}

func checkMysqlConfig(c *Config) error {
	if len(c.Mysql.Host) == 0 {
		return errors.New("empty mysql host not allowed.")
	}
	//fmt.Println(c.OutPutType)
	if strings.ToLower(c.OutPutType) == "mysql" || strings.ToLower(c.OutPutType) == "post" {
		c.OutPutType = strings.ToLower(c.OutPutType)
	} else {
		return errors.New("unKnow OutPutType.")
	}
	return nil
}
