# go-tour

go 命令行应用

## 单词格式化转换

1. 查看帮助信息

```shell
go run main.go word -h

该子命令支持各种单词格式转换，模式如下：
1：全部转大写
2：全部转小写
3：下划线转大写驼峰
4：下划线转小写驼峰
5：驼峰转下划线

Usage:
   word [flags]

Flags:
  -h, --help         help for word
  -m, --mode int8    请输入单词转换的模式
  -s, --str string   请输入单词内容

```

2. 全部转大写

```shell
go run main.go word -s=alex -m=1

# 2022/01/02 19:06:37 输出结果：ALEX
```

3. 全部转小写

```shell
go run main.go word -s=alEX -m=2

# 2022/01/02 19:08:11 输出结果：alex
```

4. 下划线转大写驼峰

```shell
go run main.go word -s=alex_pu -m=3

# 2022/01/02 19:09:13 输出结果：AlexPu
```

5. 下划线转小写驼峰

```shell
go run main.go word -s=alex_pu -m=4

# 2022/01/02 19:09:53 输出结果：alexPu
```

6. 驼峰转下划线

```shell
go run main.go word -s=AlexPu -m=5

# 2022/01/02 19:10:49 输出结果：alex_pu
```

## 便捷的时间工具

1. 查看帮助信息

```shell
go run main.go time -h

时间格式处理

Usage:
   time [flags]
   time [command]

Available Commands:
  calc        计算所需时间
  now         获取当前时间

Flags:
  -h, --help   help for time

Use " time [command] --help" for more information about a command.


------------------------------

go run main.go time calc -h

计算所需时间

Usage:
   time calc [flags]

Flags:
  -c, --calculate string    需要计算的时间，有效单位为时间戳或已格式化后的时间 
  -d, --duration string     持续时间，有效时间单位为"ns", "us" (or "µ s"), "ms", "s", "m", "h"
  -h, --help               help for calc


------------------------------

go run main.go time now -h

获取当前时间

Usage:
   time now [flags]

Flags:
  -h, --help   help for now

```

2. 获取当前时间及时间戳

```shell
go run main.go time now

# 2022/01/02 19:18:31 输出结果: 2022-01-02 19:18:31, 1641122311
```

3. 推算时间

```shell
# 查看当前时间及时间戳
go run main.go time calc -d=0s
# 2022/01/02 19:18:31 输出结果: 2022-01-02 19:18:31, 1641122311

# 推算时间（eg：当前时间减少 5 个小时）
go run main.go time calc -c="2022-01-02 19:18:31" -d=-5h
# 2022/01/02 19:43:02 输出结果: 2022-01-02 14:18:31, 1641104311
# 通过时间戳格式化时间
go run main.go time calc -c="1641104311" -d=0s
# 2022/01/02 19:50:48 输出结果: 2022-01-02 14:18:31, 1641104311

go run main.go time calc -c="2022-01-02" -d=-24h
# 2022/01/02 19:43:30 输出结果: 2022-01-01, 1640966400
go run main.go time calc -c="1640966400" -d=0h
# 2022/01/02 19:53:52 输出结果: 2022-01-01 00:00:00, 1640966400
```

## sql 语句到结构体的转换

1. 查看帮助信息

```shell
go run main.go sql struct -h

sql转换

Usage:
   sql struct [flags]

Flags:
      --charset string    请输入数据库的编码 (default "utf8mb4")
      --db string         请输入数据库名称
  -h, --help              help for struct
      --host string       请输入数据库的HOST (default "127.0.0.1:3306")
      --password string   请输入数据库的密码
      --table string      请输入表名称
      --type string       请输入数据库实例类型 (default "mysql")
      --username string   请输入数据库的账号

```

2. 数据表结构转结构体

```shell

go run main.go sql struct --username root --password 123456 --db goblog --table users

type Users struct {
         // id  bigint(20) unsigned is_nullable NO
         Id     int64   `json:"id"`
         // created_at  datetime(3) is_nullable YES
         CreatedAt      time.Time       `json:"created_at"`
         // updated_at  datetime(3) is_nullable YES
         UpdatedAt      time.Time       `json:"updated_at"`
         // name  varchar(255) is_nullable NO
         Name   string  `json:"name"`
         // email  varchar(255) is_nullable YES
         Email  string  `json:"email"`
         // password  varchar(255) is_nullable YES
         Password       string  `json:"password"`
}

func (model Users) TableName() string {
        return "users"
}


```