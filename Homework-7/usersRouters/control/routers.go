package control

import (
	"SecondTerm/Homework-7/usersRouters/middleware/cors"
	oauth2 "SecondTerm/Homework-7/usersRouters/oauth"
	"SecondTerm/Homework-7/usersRouters/oauth/jwt"
	"SecondTerm/Homework-7/usersRouters/serve"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func RoutersEntrance() {
	//同时输出到终端和日志文件
	file, _ := os.Create("ginLog.md")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	//http://121.196.155.183:8000/serve
	router := gin.Default()
	router.StaticFile("favicon.ico", "./static/favicon.ico") //加载网页图标
	router.Use(cors.Cors())                                  //跨域中间件

	oauth := router.Group("/serve/oauth")
	{
		oauth.GET("/authorize", oauth2.OAuthAuthorize) //获取授权码
		oauth.GET("/callback", oauth2.OAuthCallBack)   //获取返回的授权码
		oauth.GET("/token", oauth2.OAuthToken)         //获取token和更新token
	}

	user := router.Group("/serve/user") //用户服务
	{
		user.POST("/login", serve.Login)                    //登录
		user.POST("/register", serve.Register)              //注册
		user.PUT("/update", jwt.TokenCheck(), serve.Update) //更新个人信息
	}

	router.Run(":8001")
}
