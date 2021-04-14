package serve

import (
	"Homework-4/model"
	"Homework-4/oauth/jwt"
	ut "Homework-4/utilities"
	"errors"
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

	//检查邮箱的正确性、是否存在
	if u.Email != "" {
		if !ut.EmailCheck(u.Email) {
			resp.Return(ctx, 30002, "邮箱不合法", nil)
			return
		} else {
			ok, err := model.CheckRegisterOrNot("email", &u)
			if err != nil {
				resp.Return(ctx, 30003, "未知错误", nil)
				return
			} else if ok { //如果查找到数据
				resp.Return(ctx, 30002, "该邮箱已被注册", nil)
				return
			}

		}
	}

	//检查手机号的正确性、是否存在
	if u.Telephone != 0 {
		if !ut.TelephoneCheck(u.Telephone) {
			resp.Return(ctx, 30002, "手机号不合法", nil)
			return
		} else { //如果查找到数据
			ok, err := model.CheckRegisterOrNot("telephone", &u)
			if err != nil {
				resp.Return(ctx, 30003, "未知错误", nil)
				return
			} else if ok {
				resp.Return(ctx, 30002, "该手机号已被注册", nil)
				return
			}
		}
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

//更新用户数据
func Update(ctx *gin.Context) {
	var (
		resp ut.Resp
		err  error
	)

	tem, _ := strconv.ParseUint(ctx.Query("user_id"), 10, 64) //用户id
	userId := uint(tem)                                       //uint64转uint
	operateType, _ := strconv.Atoi(ctx.Query("type"))         //操作类型
	content := ctx.Query("content")                           //更改成什么内容

	//简单排除类型错误
	if operateType > 4 || operateType < 1 {
		resp.Return(ctx, 12001, "类型不合法", nil)
		return
	}

	if operateType == 1 { //更新手机号
		newTele, _ := strconv.ParseInt(content, 10, 64)
		if ut.TelephoneCheck(newTele) {
			err = model.ChangeUserTelephone(userId, newTele)
		}
	} else if operateType == 2 { //更新邮箱
		if ut.EmailCheck(content) {
			err = model.ChangeUserEmail(userId, content)
		}
	} else if operateType == 3 { //更新昵称
		if ut.UserNameCheck(content) {
			err = model.ChangeUserNickname(userId, content)
		}
	} else if operateType == 4 { //更新签名
		err = model.ChangeUserSignature(userId, content)
	}

	if err != nil {
		if errors.Is(err, ut.ErrorUserNotExist) {
			resp.Return(ctx, 12003, "用户数据错误", nil)
		} else {
			resp.Return(ctx, 12002, "未知错误", nil)
		}
		return
	}

	resp.Return(ctx, 1200, "更新成功", nil)
}
