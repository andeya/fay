// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/henrylee2cn/think/model"
	"github.com/henrylee2cn/thinkgo"
)

var appname string
var crupath, _ = os.Getwd()

func main() {
	thinkgo.Printf("os.Args: %#v", os.Args)
	if len(os.Args) < 2 {
		help()
		return
	}
	switch os.Args[1] {
	case "new":
		newapp(os.Args[2:])
	}
}
func help() {

}
func newappHelp() {

}

func cleanCrupath() {
	var err error
	crupath = strings.TrimSpace(crupath)
	crupath, err = filepath.Abs(crupath)
	if err != nil {
		thinkgo.Fatalf("[think] Create project fail: %s", err)
	}
	crupath = strings.Replace(crupath, `\`, `/`, -1)
	crupath = strings.TrimRight(crupath, "/") + "/"
}

func newapp(args []string) {
	switch len(args) {
	case 1:
		appname = args[0]
		crupath = filepath.Join(crupath, appname)
	case 2:
		appname = args[0]
		crupath = args[1]
	default:
		newappHelp()
		return
	}
	cleanCrupath()
	thinkgo.Printf("[think] Create a thinkgo project named `%s` in the `%s` path.", appname, crupath)
	if isExist(crupath) {
		thinkgo.Printf("[think] The project path has conflic, do you want to build in: %s\n", crupath)
		thinkgo.Printf("[think] Do you want to overwrite it? [yes|no]]  ")
		if !askForConfirmation() {
			thinkgo.Fatalf("[think] cancel...")
			return
		}
	}

	exit := make(chan bool)
	thinkgo.Printf("[think] Start create project...")

	model.SimplePro(crupath, appname)

	thinkgo.Printf("[think] Create was successful")

	if err := os.Chdir(crupath); err != nil {
		thinkgo.Fatalf("[think] Create project fail: %v", err)
	}
	autobuild()
	newWatcher()
	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}
