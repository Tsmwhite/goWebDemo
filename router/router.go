package router

import (
	"github.com/gin-gonic/gin"
	"z3/controllers"
)

func RegRouter(router *gin.Engine)  {
	//访问路由不存在时
	router.NoRoute(controllers.Handle404)

	//登陆注册
	router.POST("/login",controllers.LoginAuth)
	router.POST("/register",controllers.Register)

	//个人中心
	userCenter := router.Group("/user",controllers.AuthVerify())
	{
		userCenter.GET("/info")
		userCenter.POST("/headimg",controllers.ChangeHeadimg)
		userCenter.POST("/password",controllers.ChangePassword)
		userCenter.POST("/changeinfo")
	}
}