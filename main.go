package main

import (
	"github.com/zhutingle/gotrix/server"
)

func main() {

	server.GotrixServer()

	//signs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)
	//
	//signal.Notify(signs, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP, syscall.SIGTERM)
	//
	//fmt.Println(os.Getpid())
	//
	//go func() {
	//	for {
	//		sig := <-signs
	//		fmt.Println(sig)
	//		done <- true
	//	}
	//
	//}()
	//
	//for i := 0; i < 10; i++ {
	//	<-done
	//	fmt.Println("Done Signal")
	//}

}
