# fay    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/fay)    ![Fay goreportcard](https://goreportcard.com/badge/github.com/henrylee2cn/fay)

Fay is the deployment tool of faygo web framework.

Faygo is a fast and concise Go Web framework that can be used to develop high-performance web app(especially API) with fewer codes. Just define a struct Handler, Faygo will automatically bind/verify the request parameters and generate the online API doc. [Go to faygo](https://github.com/henrylee2cn/faygo)

[简体中文](https://github.com/henrylee2cn/fay/blob/master/README_ZH.md)

## Features

- Create, compile and run (monitor changes) a new faygo project
- Compile and run (monitor changes) an any existing go project
- Provides a meta-programming toolkit for faygo

## Requirements

Go Version ≥1.8

## Download and install

```sh
go get -u -v github.com/henrylee2cn/fay
```

## Usage

```
        fay command [arguments]

The commands are:
        new        create, compile and run (monitor changes) a new faygo project
        run        compile and run (monitor changes) an any existing go project

fay new appname [apptpl]
        appname    specifies the path of the new faygo project
        apptpl     optionally, specifies the faygo project template type

fay run [appname]
        appname    optionally, specifies the path of the new project
```