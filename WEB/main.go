package main

import (
	"github.com/zhutingle/gotrix/global"
	"github.com/zhutingle/gotrix/server"
)

func main() {

	// exec.Command("/bin/stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// exec.Command("/bin/stty", "-F", "/dev/tty", "/echo").Run()
	// defer exec.Command("/bin/stty", "-F", "/dev/tty", "echo").Run()

	//exec.Command("ls","/l").Output()
	//
	//
	//var password string
	//fmt.Scanln(&password)
	//fmt.Println(password)

	global.Config.Args.ConfigFile = "/Users/zhutingle/Documents/workspace/gotrix/src/github.com/zhutingle/gotrix/WEB/gotrix.conf"

	server.GotrixServer()

}
