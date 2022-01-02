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