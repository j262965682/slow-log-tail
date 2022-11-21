# slow-log-tail
实时拉取mysql慢查询输出到数据库

# 原理
通过对slow-log文件的实时增量分析，产生慢查询数据，可以输出配置 插入数据库或者通过url post传出去

### 配置文件
```yaml
# mysql or post
outPutType: 'post'
mysql:
  host: '172.0.0.1'
  port: 3306
  user: 'slow'
  password: 'slow'
  dbname: 'onecool'
  tablename: 't_slow_log'
  other: 'charset=utf8mb4&parseTime=True&loc=Local'
post:
  host: '172.0.0.1'
  url:  '172.0.0.1'
slowLog:
  path: 'D:\GitHubPro\22222\slow-log-tail\slowlog\slow-log.log'
  instance: '测试'
  ignoreuser: ['root']
  longquerytime: 0.4

```

### 编译
```shell
set GOOS=linux
set GOARCH=amd64
go build -o "slow-log-tail"
```

### 执行
```shell
chmod +x slow-log-tail
nohup ./slow-log-tail >/dev/null 2>slow-log-nohup.log & 
```