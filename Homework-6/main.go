package main

import (
	"SecondTerm/Homework-6/sky"
)

func main() {
	router := sky.Default()

	router.GET("/hello/:name", func(ctx *sky.Context) {
		ctx.String(200, ctx.Param("name"))
	})

	stop := router.Group("/stop")
	stop.GET("/*path", func(ctx *sky.Context) {
		ctx.JSON(200, sky.H{
			"info": ctx.Param("path"),
		})
	})

	router.Run()
}

//func main() {
//	router := gin.Default()
//	router.GET("", func(context *gin.Context) {
//		context.Next()
//		context.Cookie()
//	})
//	router.Run()
//}
