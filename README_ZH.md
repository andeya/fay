# think    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/think)    ![Think goreportcard](https://goreportcard.com/badge/github.com/henrylee2cn/think)

Thinkgo WEB框架的开发工具。

[Go to thinkgo](https://github.com/henrylee2cn/thinkgo)

## 特点

- 新建、编译、运行（实时监控文件变动）一个新的thinkgo项目
- 支持运行任意的golang程序
- 提供Thinkgo的元编程工具包

## 下载并安装

```sh
go get -u -v github.com/henrylee2cn/think
```

## 用法

```
        think command [arguments]

The commands are:
        new        创建、编译和运行（监控文件变化）一个新的thinkgo项目
        run        编译和运行（监控文件变化）任意一个已存在的golang项目

think new appname [apptpl]
        appname    指定新thinkgo项目的创建目录
        apptpl     指定一个thinkgo项目模板（可选）

think run [appname]
        appname    指定待运行的golang项目路径（可选）
```
