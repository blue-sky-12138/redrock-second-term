[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] GET    /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (3 handlers)
[GIN-debug] HEAD   /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (3 handlers)
[GIN-debug] GET    /serve/oauth/authorize    --> Homework-4/oauth.OAuthAuthorize (4 handlers)
[GIN-debug] GET    /serve/oauth/callback     --> Homework-4/oauth.OAuthCallBack (4 handlers)
[GIN-debug] GET    /serve/oauth/token        --> Homework-4/oauth.OAuthToken (4 handlers)
[GIN-debug] POST   /serve/user/login         --> Homework-4/serve.Login (4 handlers)
[GIN-debug] POST   /serve/user/register      --> Homework-4/serve.Register (4 handlers)
[GIN-debug] PUT    /serve/user/update        --> Homework-4/serve.Update (5 handlers)
[GIN-debug] GET    /serve/video/comment      --> Homework-4/serve.GetVideoComments (4 handlers)
[GIN-debug] GET    /serve/video/information  --> Homework-4/serve.GetVideoInformation (4 handlers)
[GIN-debug] GET    /serve/video/barrage      --> Homework-4/serve.GetVideoBarrages (4 handlers)
[GIN-debug] GET    /serve/video/path         --> Homework-4/serve.GetVideoPath (4 handlers)
[GIN-debug] PUT    /serve/video/operation    --> Homework-4/serve.OperateVideo (5 handlers)
[GIN-debug] POST   /serve/video/comment      --> Homework-4/serve.AddComment (5 handlers)
[GIN-debug] GET    /serve/download/user/head/:id/:fileName --> Homework-4/serve.GetUserHead (4 handlers)
[GIN-debug] GET    /serve/download/video/cover/:bvCode/:fileName --> Homework-4/serve.GetVideoFile (4 handlers)
[GIN-debug] GET    /serve/download/video/file/:bvCode/:fileName --> Homework-4/serve.GetVideoFile (4 handlers)
[GIN-debug] PUT    /serve/upload/user/head   --> Homework-4/serve.UpdateUserHead (5 handlers)
[GIN-debug] Listening and serving HTTP on :8000
[GIN] 2021/04/14 - 22:59:05 | 301 |    415.9226ms |             ::1 | GET      "/serve/oauth/authorize?response_type=code&client_id=123456&scope=read&redirect_url=http://localhost:8000/serve/oauth/callback&state="
[GIN] 2021/04/14 - 22:59:05 | 200 |            0s |             ::1 | GET      "/serve/oauth/callback?code=5881967194749572270"
[GIN] 2021/04/14 - 22:59:25 | 200 |    414.5577ms |             ::1 | GET      "/serve/oauth/authorize?response_type=code&client_id=123456&scope=read&state="
[GIN] 2021/04/14 - 22:59:32 | 200 |    433.1565ms |             ::1 | GET      "/serve/oauth/token?grant_type=authorization_code&client_id=123456&scope=read&client_secret=RedRock&code=8798430169309599164"
