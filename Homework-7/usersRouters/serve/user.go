package serve

import (
	"SecondTerm/Homework-7/usersRouters/oauth/jwt"
	"SecondTerm/Homework-7/usersRouters/usersCentercall"
	ut "SecondTerm/Homework-7/usersRouters/utilities"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//用户登录
func Login(ctx *gin.Context) {
	var (
		login ut.UserInformation
		check ut.UserInformation
		resp  ut.Resp
	)
	err := ctx.ShouldBind(&login)
	if err != nil {
		resp.Return(ctx, 20003, "未知错误", nil)
		ut.LogError("GetLoginInformation Error", err)
		return
	}

	code, message, err := usersCentercall.Login(login.Name, login.Password)
	if err != nil {
		resp.Return(ctx, 20003, "未知错误", nil)
		ut.LogError("CallLogin Error", err)
		return
	}

	if code != 200 {
		resp.Return(ctx, code, message, nil)
	} else {
		fmt.Println(check.ID, check.Name)
		cookie := &http.Cookie{
			Name:     "user",
			Value:    jwt.NewJWT(check.ID, check.Name).Token,
			MaxAge:   100000,
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(ctx.Writer, cookie)

		resp.Return(ctx, 200, "登陆成功", nil)
	}
}

//用户注册
func Register(ctx *gin.Context) {
	var (
		u    ut.UserInformation //用于接收传输的数据
		resp ut.Resp            //响应体
	)
	err := ctx.ShouldBind(&u)
	if err != nil {
		ut.LogError("GetRegisterInformation Error", err)
		resp.Return(ctx, 30003, "未知错误", nil)
		return
	}

	resp.Code, resp.Message, err = usersCentercall.Register(u.Name, u.Password, u.Email, u.Telephone)
	if err != nil {
		ut.LogError("CallRegister Error", err)
		resp.Return(ctx, 30003, "未知错误", nil)
		return
	}

	resp.Return(ctx, resp.Code, resp.Message, nil)
}

//更新用户数据
func Update(ctx *gin.Context) {
	var (
		resp ut.Resp
		err  error
	)

	userid, _ := strconv.Atoi(ctx.Query("user_id"))   //用户id
	operateType, _ := strconv.Atoi(ctx.Query("type")) //操作类型
	content := ctx.Query("content")                   //更改成什么内容

	resp.Code, resp.Message, err = usersCentercall.Update(userid, operateType, content)
	if err != nil {
		ut.LogError("CallUpdate Error", err)
		resp.Return(ctx, 30003, "未知错误", nil)
	}

	resp.Return(ctx, resp.Code, resp.Message, nil)
}
