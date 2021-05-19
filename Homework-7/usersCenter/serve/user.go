package serve

import (
	"SecondTerm/Homework-7/usersCenter/model"
	"SecondTerm/Homework-7/usersCenter/oauth/jwt"
	ut "SecondTerm/Homework-7/usersCenter/utilities"
	"errors"
	"github.com/MashiroC/begonia"
	"strconv"
)

func Entrance() {
	s := begonia.NewServer()
	s.Register("UsersServer", &RPCUserServer{})
	for {
		s.Wait()
	}
}

type RPCUserServer struct{}

func (r *RPCUserServer) Login(name string, password string) (code int, message string, data string) {
	var (
		login ut.UserInformation
		check ut.UserInformation
	)
	login.Name = name
	login.Password = password

	check.Name = login.Name
	ok, err := model.GetUserPassword(&check)
	if err != nil {
		return 20003, "未知错误", ""
	}

	if !ok { //检查用户名是否存在
		return 20001, "该用户不存在", ""
	} else if ut.Cryptography(login.Password, check.Md5salt) != check.Password { //检查密码是否正确
		return 20002, "密码错误", ""
	} else {
		return 200, "登陆成功", jwt.NewJWT(check.ID, check.Name).Token
	}
}

//用户注册
func (r *RPCUserServer) Register(name string, password string, email string, telephone int64) (code int, message string) {
	u := ut.UserInformation{
		Name:      name,
		Password:  password,
		Telephone: telephone,
		Email:     email,
	}

	//检查邮箱的正确性、是否存在
	if u.Email != "" {
		if !ut.EmailCheck(u.Email) {
			return 30002, "邮箱不合法"
		} else {
			ok, err := model.CheckRegisterOrNot("email", &u)
			if err != nil {
				return 30003, "未知错误"
			} else if ok { //如果查找到数据
				return 30002, "该邮箱已被注册"
			}

		}
	}

	//检查手机号的正确性、是否存在
	if u.Telephone != 0 {
		if !ut.TelephoneCheck(u.Telephone) {
			return 30002, "手机号不合法"
		} else { //如果查找到数据
			ok, err := model.CheckRegisterOrNot("telephone", &u)
			if err != nil {
				return 30003, "未知错误"
			} else if ok {
				return 30002, "该手机号已被注册"
			}
		}
	}

	//检查用户名的正确性、是否存在
	if u.Name != "" {
		if !ut.UserNameCheck(u.Name) {
			return 30002, "用户名不符合要求"
		} else { //如果查找到数据
			ok, err := model.CheckRegisterOrNot("name", &u)
			if err != nil {
				return 30003, "未知错误"
			} else if ok {
				return 30002, "该用户名已被注册"
			}
		}
	}

	//检查密码规范性
	if !ut.PasswordCheck(u.Password) {
		return 30002, "密码不规范"
	}

	//md5加盐加密
	u.Password, u.Md5salt = ut.CryptographyNow(u.Password)

	err := model.PostRegisterInformation(&u)
	if err != nil {
		return 30003, "未知错误"
	}

	return 300, "注册成功"
}

//更新用户数据
func (r *RPCUserServer) Update(userId int, operateType int, content string) (code int, message string) {
	var (
		err     error
		usersId = uint(userId)
	)

	//简单排除类型错误
	if operateType > 4 || operateType < 1 {
		return 12001, "类型不合法"
	}

	if operateType == 1 { //更新手机号
		newTele, _ := strconv.ParseInt(content, 10, 64)
		if ut.TelephoneCheck(newTele) {
			err = model.ChangeUserTelephone(usersId, newTele)
		}
	} else if operateType == 2 { //更新邮箱
		if ut.EmailCheck(content) {
			err = model.ChangeUserEmail(usersId, content)
		}
	} else if operateType == 3 { //更新昵称
		if ut.UserNameCheck(content) {
			err = model.ChangeUserNickname(usersId, content)
		}
	} else if operateType == 4 { //更新签名
		err = model.ChangeUserSignature(usersId, content)
	}

	if err != nil {
		if errors.Is(err, ut.ErrorUserNotExist) {
			return 12003, "用户数据错误"
		} else {
			return 12002, "未知错误"
		}
		return
	}

	return 1200, "更新成功"
}
