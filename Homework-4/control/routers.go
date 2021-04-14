package control

import (
	"Homework-4/middleware/cors"
	oauth2 "Homework-4/oauth"
	"Homework-4/oauth/jwt"
	"Homework-4/serve"
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
	router.Static("static", "./static/")                     //加载静态文件夹
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

	video := router.Group("/serve/video") //视频服务
	{
		video.GET("/comment", serve.GetVideoComments)                 //获取视频评论
		video.GET("/information", serve.GetVideoInformation)          //获取视频的元数据
		video.GET("/barrage", serve.GetVideoBarrages)                 //获取视频弹幕
		video.GET("/path", serve.GetVideoPath)                        //获取视频地址
		video.PUT("/operation", jwt.TokenCheck(), serve.OperateVideo) //用户对视频进行点赞等操作
		video.POST("/comment", jwt.TokenCheck(), serve.AddComment)    //添加评论
	}

	download := router.Group("/serve/download")
	{
		download.GET("/user/head/:id/:fileName", serve.GetUserHead)        //获取头像
		download.GET("/video/cover/:bvCode/:fileName", serve.GetVideoFile) //获取视频封面
		download.GET("/video/file/:bvCode/:fileName", serve.GetVideoFile)  //获取视频文件本体
	}

	upload := router.Group("/serve/upload")
	{
		upload.PUT("/user/head", jwt.TokenCheck(), serve.UpdateUserHead) //更新用户头像
		//upload.POST("/video/file_one", serve.UploadVideoOne)   //上传单个视频(投稿)
		//upload.POST("/video/file_more", serve.UploadVideoMore) //上传多个视频(投稿)
	}

	router.Run(":8000")
}
