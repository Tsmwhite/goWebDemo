package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"z3/webscoket"
	"z3/router"
	"z3/common/mail"
)

func main (){
	res := mail.ResMail().SetMessage("577689878@qq.com","test-send-mail","<h1>go go go </h1>").Send()
	fmt.Println(res)
	return
	go func(){
		//监听websocket
		go webscoket.Manager.ListenEvent()
		http.HandleFunc("/ws", webscoket.WsHandler)
		http.ListenAndServe(":8888",nil)
	}()

	//注册gin路由 开启端口监听
	ginEngine := gin.Default()
	router.RegRouter(ginEngine)
	ginEngine.Run(":8080")
}