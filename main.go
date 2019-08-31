package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"z3/controllers"
	"z3/webscoket"
)

func main (){
	go func(){
		//监听websocket
		go webscoket.Manager.ListenEvent()
		http.HandleFunc("/ws", webscoket.WsHandler)
		http.ListenAndServe(":8888",nil)
	}()

	router := gin.Default()
	router.NoRoute(Handle404)
	router.POST("/login",controllers.LoginAuth)
	router.POST("/register",controllers.Register)

	userCenter := router.Group("/user")
	userCenter.Use(controllers.AuthVerify())
	{
		userCenter.GET("/info")
		userCenter.POST("/headimg",controllers.ChangeHeadimg)
		userCenter.POST("/password",controllers.ChangePassword)
		userCenter.POST("/changeinfo")
	}

	router.Run(":8080")
}

func Handle404(c *gin.Context) {
	controllers.ResError("Not Fund 404",c)
}
