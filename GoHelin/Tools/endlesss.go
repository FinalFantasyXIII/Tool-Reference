package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"syscall"
)

func main(){
	//注册一个http服务
	router := gin.Default()
	router.GET("/helin/friday", func(context *gin.Context) {
		context.JSON(0,"hello")
	})

	server := endless.NewServer(":8080",router)

	//注册http服务结束前需要处理的工作
	server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT] = append(
		server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT],
		atExit)

	//绑定启动,只要从此处返回，err必定不为nil
	err := server.ListenAndServe()
	if err != nil{
		fmt.Println(err)
	}

	//执行后续工作
	fmt.Println("over ...")
}


func atExit() {
	log.Println("exit ...")
}



