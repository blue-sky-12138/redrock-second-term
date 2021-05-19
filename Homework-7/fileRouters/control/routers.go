package control

import (
	"SecondTerm/Homework-7/fileRouters/middleware/cors"
	"SecondTerm/Homework-7/fileRouters/oauth/jwt"
	"SecondTerm/Homework-7/fileRouters/serve"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func RoutersEntrance() {
	//同时输出到终端和日志文件
	file, _ := os.Create("ginLog.md")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	//http://121.196.155.183:8003/serve
	router := gin.Default()
	router.Static("static", "./static/")                     //加载静态文件夹
	router.StaticFile("favicon.ico", "./static/favicon.ico") //加载网页图标
	router.Use(cors.Cors())                                  //跨域中间件

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

	router.Run(":8001")
}
