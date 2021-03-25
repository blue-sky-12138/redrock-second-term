[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] GET    /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (3 handlers)
[GIN-debug] HEAD   /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (3 handlers)
[GIN-debug] POST   /serve/user/login         --> Homework-2/serve.Login (4 handlers)
[GIN-debug] POST   /serve/user/register      --> Homework-2/serve.Register (4 handlers)
[GIN-debug] PUT    /serve/user/update        --> Homework-2/serve.Update (4 handlers)
[GIN-debug] GET    /serve/video/comment      --> Homework-2/serve.GetVideoComments (4 handlers)
[GIN-debug] GET    /serve/video/information  --> Homework-2/serve.GetVideoInformation (4 handlers)
[GIN-debug] GET    /serve/video/barrage      --> Homework-2/serve.GetVideoBarrages (4 handlers)
[GIN-debug] GET    /serve/video/path         --> Homework-2/serve.GetVideoPath (4 handlers)
[GIN-debug] PUT    /serve/video/operation    --> Homework-2/serve.OperateVideo (4 handlers)
[GIN-debug] POST   /serve/video/comment      --> Homework-2/serve.AddComment (4 handlers)
[GIN-debug] GET    /serve/download/user/head/:id/:fileName --> Homework-2/serve.GetUserHead (4 handlers)
[GIN-debug] GET    /serve/download/video/cover/:bvCode/:fileName --> Homework-2/serve.GetVideoFile (4 handlers)
[GIN-debug] GET    /serve/download/video/file/:bvCode/:fileName --> Homework-2/serve.GetVideoFile (4 handlers)
[GIN-debug] PUT    /serve/upload/user/head   --> Homework-2/serve.UpdateUserHead (4 handlers)
[GIN-debug] Listening and serving HTTP on :8000
[GIN] 2021/03/25 - 09:49:52 | 200 |     1.595205s |       127.0.0.1 | GET      "/serve/video/comment?bv_code=BV1No4y1d7tA"
