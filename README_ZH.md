# fay    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/fay)    ![Fay goreportcard](https://goreportcard.com/badge/github.com/henrylee2cn/fay)

Fay 是Go Web框架 Faygo 的开发工具。

Faygo 是一款快速、简洁的Go Web框架，可用极少的代码开发出高性能的Web应用程序（尤其是API接口）。只需定义 struct Handler，Faygo 就能自动绑定、验证请求参数并生成在线API文档。

[Go to faygo](https://github.com/henrylee2cn/faygo)

## 特点

- 新建、编译、运行（实时监控文件变动）一个新的faygo项目
- 支持运行任意的golang程序
- 提供Faygo的元编程工具包


## 安装要求

Go Version ≥1.8

## 下载并安装

```sh
go get -u -v github.com/henrylee2cn/fay
```

## 用法

```
        fay command [arguments]

The commands are:
        new        创建、编译和运行（监控文件变化）一个新的faygo项目
        run        编译和运行（监控文件变化）任意一个已存在的golang项目

fay new appname [apptpl]
        appname    指定新faygo项目的创建目录
        apptpl     指定一个faygo项目模板（可选）

fay run [appname]
        appname    指定待运行的golang项目路径（可选）
```
