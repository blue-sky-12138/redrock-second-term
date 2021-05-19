package control

import (
	"SecondTerm/Homework-7/videosRouters/middleware/cors"
	"SecondTerm/Homework-7/videosRouters/oauth/jwt"
	"SecondTerm/Homework-7/videosRouters/serve"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func RoutersEntrance() {
	//同时输出到终端和日志文件
	file, _ := os.Create("ginLog.md")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	//http://121.196.155.183:8002/serve
	router := gin.Default()
	router.StaticFile("favicon.ico", "./static/favicon.ico") //加载网页图标
	router.Use(cors.Cors())                                  //跨域中间件

	video := router.Group("/serve/video") //视频服务
	{
		video.GET("/comment", serve.GetVideoComments)                 //获取视频评论
		video.GET("/information", serve.GetVideoInformation)          //获取视频的元数据
		video.GET("/barrage", serve.GetVideoBarrages)                 //获取视频弹幕
		video.GET("/path", serve.GetVideoPath)                        //获取视频地址
		video.PUT("/operation", jwt.TokenCheck(), serve.OperateVideo) //用户对视频进行点赞等操作
		video.POST("/comment", jwt.TokenCheck(), serve.AddComment)    //添加评论
	}

	router.Run(":8002")
}
