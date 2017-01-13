# think    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/think)    ![Think goreportcard](https://goreportcard.com/badge/github.com/henrylee2cn/think)

The deployment tools of thinkgo web frameware.

[Go to thinkgo](https://github.com/henrylee2cn/thinkgo)

[简体中文](https://github.com/henrylee2cn/think/blob/master/README_ZH.md)

## Features

- Create, compile and run (monitor changes) a new thinkgo project
- Compile and run (monitor changes) an any existing go project
- Provides a meta-programming toolkit for thinkgo

## Download and install

```sh
go get -u -v github.com/henrylee2cn/think
```

## Usage

```
        think command [arguments]

The commands are:
        new        create, compile and run (monitor changes) a new thinkgo project
        run        compile and run (monitor changes) an any existing go project

think new appname [apptpl]
        appname    specifies the path of the new thinkgo project
        apptpl     optionally, specifies the thinkgo project template type

think run [appname]
        appname    optionally, specifies the path of the new project
```