[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] GET    /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (3 handlers)
[GIN-debug] HEAD   /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (3 handlers)
[GIN-debug] GET    /serve/download/user/head/:id/:fileName --> SecondTerm/Homework-7/fileRouters/serve.GetUserHead (4 handlers)
[GIN-debug] GET    /serve/download/video/cover/:bvCode/:fileName --> SecondTerm/Homework-7/fileRouters/serve.GetVideoFile (4 handlers)
[GIN-debug] GET    /serve/download/video/file/:bvCode/:fileName --> SecondTerm/Homework-7/fileRouters/serve.GetVideoFile (4 handlers)
[GIN-debug] PUT    /serve/upload/user/head   --> SecondTerm/Homework-7/fileRouters/serve.UpdateUserHead (5 handlers)
[GIN-debug] Listening and serving HTTP on :8001
