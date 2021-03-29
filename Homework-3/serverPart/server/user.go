package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"serverPart/model"
	ut "serverPart/utilities"
)

//用户登录
func Login(ctx *gin.Context) {
	var (
		login ut.UserInformation
		check ut.UserInformation
		resp  ut.Resp
	)
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		resp.Return(ctx, 20003, "未知错误", nil)
		log.Println("GetLoginInformation Error", err)
		return
	}

	check.Name = login.Name
	ok, err := model.GetUserPassword(&check)
	if err != nil {
		resp.Return(ctx, 20003, "未知错误", nil)
		return
	}

	if !ok { //检查用户名是否存在
		resp.Return(ctx, 20001, "该用户不存在", nil)
	} else if ut.Cryptography(login.Password, check.Md5salt) != check.Password { //检查密码是否正确
		resp.Return(ctx, 20002, "密码错误", nil)
	} else {
		resp.Return(ctx, 200, "登陆成功", nil)
	}
}

//用户注册
func Register(ctx *gin.Context) {
	var (
		u    ut.UserInformation //用于接收传输的数据
		resp ut.Resp            //响应体
	)
	err := ctx.ShouldBindJSON(&u)
	if err != nil {
		log.Println("GetRegisterInformation Error", err)
		resp.Return(ctx, 30003, "未知错误", nil)
		return
	}

	//检查用户名的正确性、是否存在
	if u.Name != "" {
		if !ut.UserNameCheck(u.Name) {
			resp.Return(ctx, 30002, "用户名不符合要求", nil)
			return
		} else { //如果查找到数据
			ok, err := model.CheckRegisterOrNot("name", &u)
			if err != nil {
				resp.Return(ctx, 30003, "未知错误", nil)
				return
			} else if ok {
				resp.Return(ctx, 30002, "该用户名已被注册", nil)
				return
			}
		}
	}

	//检查密码规范性
	if !ut.PasswordCheck(u.Password) {
		resp.Return(ctx, 30002, "密码不规范", nil)
		return
	}

	//md5加盐加密
	u.Password, u.Md5salt = ut.CryptographyNow(u.Password)

	err = model.PostRegisterInformation(&u)
	if err != nil {
		resp.Return(ctx, 30003, "未知错误", nil)
		return
	}

	resp.Return(ctx, 300, "注册成功", nil)
}
