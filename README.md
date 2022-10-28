# go-tour

go 命令行应用

## 安装

```shell
# Go version >= 1.16
$ go install github.com/pudongping/go-tour@latest

# Go version < 1.16
$ go get -u github.com/pudongping/go-tour
```

## 使用

> 更多命令，请查看 `源码运行` 示例

```shell
# 查看帮助信息
go-tour -h

# 示例：
# 比如，打印 sql 结构体。
go-tour sql struct --host 127.0.0.1:3306 --username root --password 123456 --db goblog --table users
```

## 源码运行

### 单词格式化转换

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

### 时间处理工具

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

### sql 语句到结构体的转换（支持 gorm 和 xorm）

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
      --host string       请输入数据库的 HOST (default "127.0.0.1:3306")
      --orm string        请输入你想输出的 ORM 模型名称，支持：gorm 或 xorm (default "gorm")
      --password string   请输入数据库的密码
      --table string      请输入表名称
      --type string       请输入数据库实例类型 (default "mysql")
      --username string   请输入数据库的账号
```

2. 数据表结构转结构体

默认为 `gorm`

```shell

go run main.go sql struct --host 127.0.0.1:3306 --username root --password 123456 --db goblog --table users

type Users struct {
	// user_id  auto_increment PRI bigint(20) unsigned is_nullable: NO
	UserId int64 `gorm:"column:user_id;primaryKey;autoIncrement;not null;" json:"user_id"`
	// 用户名    varchar(255) is_nullable: YES
	Name string `gorm:"column:name;" json:"name"`
	// 邮箱   UNI varchar(80) is_nullable: YES
	Email string `gorm:"column:email;unique;" json:"email"`
	// 手机号码   UNI varchar(40) is_nullable: YES
	Phone string `gorm:"column:phone;unique;" json:"phone"`
	// 自定义头像    varchar(255) is_nullable: NO
	Avatar string `gorm:"column:avatar;not null;" json:"avatar"`
	// 账号   UNI varchar(40) is_nullable: YES
	Account string `gorm:"column:account;unique;" json:"account"`
	// 密码    varchar(255) is_nullable: NO
	Password string `gorm:"column:password;not null;" json:"password"`
	// 状态，1=启用，2=禁用    tinyint(4) is_nullable: NO  default_value: 1
	Status int8 `gorm:"column:status;not null;" json:"status"`
	// 状态，1=是管理员，2=不是管理员    tinyint(4) is_nullable: NO  default_value: 1
	IsAdmin int8 `gorm:"column:is_admin;not null;" json:"is_admin"`
	// 简介    varchar(255) is_nullable: NO
	Introduction string `gorm:"column:introduction;not null;" json:"introduction"`
	// 登录时间    int(11) unsigned is_nullable: NO  default_value: 0
	LoginAt int `gorm:"column:login_at;not null;" json:"login_at"`
	// 创建时间    int(11) unsigned is_nullable: NO  default_value: 0
	CreatedAt int `gorm:"column:created_at;not null;" json:"created_at"`
	// 更新时间    int(11) unsigned is_nullable: NO  default_value: 0
	UpdatedAt int `gorm:"column:updated_at;not null;" json:"updated_at"`
	// 软删除时间    int(11) unsigned is_nullable: NO  default_value: 0
	DeletedAt int `gorm:"column:deleted_at;not null;" json:"deleted_at"`
}

func (model Users) TableName() string {
	return "users"
}

```

如若需要转换为 `xorm` 时

```shell

go run main.go sql struct --username root --password 123456 --db card_robot --table users --orm xorm

type Users struct {
	// user_id  auto_increment PRI bigint(20) unsigned is_nullable: NO
	UserId int64 `xorm:"bigint(20) pk autoincr notnull 'user_id'" json:"user_id"`
	// 用户名    varchar(255) is_nullable: YES
	Name string `xorm:"varchar(255) 'name'" json:"name"`
	// 邮箱   UNI varchar(80) is_nullable: YES
	Email string `xorm:"varchar(80) unique 'email'" json:"email"`
	// 手机号码   UNI varchar(40) is_nullable: YES
	Phone string `xorm:"varchar(40) unique 'phone'" json:"phone"`
	// 自定义头像    varchar(255) is_nullable: NO
	Avatar string `xorm:"varchar(255) notnull 'avatar'" json:"avatar"`
	// 账号   UNI varchar(40) is_nullable: YES
	Account string `xorm:"varchar(40) unique 'account'" json:"account"`
	// 密码    varchar(255) is_nullable: NO
	Password string `xorm:"varchar(255) notnull 'password'" json:"password"`
	// 状态，1=启用，2=禁用    tinyint(4) is_nullable: NO  default_value: 1
	Status int8 `xorm:"tinyint(4) notnull 'status'" json:"status"`
	// 状态，1=是管理员，2=不是管理员    tinyint(4) is_nullable: NO  default_value: 1
	IsAdmin int8 `xorm:"tinyint(4) notnull 'is_admin'" json:"is_admin"`
	// 简介    varchar(255) is_nullable: NO
	Introduction string `xorm:"varchar(255) notnull 'introduction'" json:"introduction"`
	// 登录时间    int(11) unsigned is_nullable: NO  default_value: 0
	LoginAt int `xorm:"int(11) notnull 'login_at'" json:"login_at"`
	// 创建时间    int(11) unsigned is_nullable: NO  default_value: 0
	CreatedAt int `xorm:"int(11) notnull created 'created_at'" json:"created_at"`
	// 更新时间    int(11) unsigned is_nullable: NO  default_value: 0
	UpdatedAt int `xorm:"int(11) notnull updated 'updated_at'" json:"updated_at"`
	// 软删除时间    int(11) unsigned is_nullable: NO  default_value: 0
	DeletedAt int `xorm:"int(11) notnull deleted 'deleted_at'" json:"deleted_at"`
}

func (model Users) TableName() string {
	return "users"
}

```


### json 字符串转 struct 结构体

1. 查看帮助信息

```shell

go run main.go json struct -h 

json转换

Usage:
   json struct [flags]

Flags:
  -h, --help         help for struct
  -s, --str string   请输入json字符串

```

2. json 字符串转结构体

```shell

go run main.go json struct -s='{"code":200,"msg":"Success","data":{"name":"alex","age":26,"favorite":["basketball","pingpong-ball"],"lucky-number":[6, 8]}}'

2022/01/03 04:20:09 输出结果: 

type Tour struct {
    Code float64
    Msg string
    Data map[string]interface {}
}

```