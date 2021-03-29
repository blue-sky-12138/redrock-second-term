package model

import (
	"regexp"
	"serverPart/dao"
	ut "serverPart/utilities"
	"strconv"
	"strings"
)

//获取用户密码。
func GetUserPassword(u *ut.UserInformation) (bool, error) {
	var (
		err  error
		from string //以什么作为搜索目标
	)
	if strings.Contains(u.Name, "@") { //判断是否有@，有即为邮箱登录
		u.Email, u.Name = u.Name, ""
		from = "email"
	} else {
		reg := regexp.MustCompile("[^0-9]") //判断是否为纯数字，是即为手机号登录
		if reg.MatchString(u.Name) {
			from = "name"
		} else {
			u.Telephone, _ = strconv.ParseInt(u.Name, 10, 64)
			u.Name = ""
			from = "telephone"
		}
	}

	err = dao.DB.Select("id", "password", "md5salt").Where(u, from).Find(u).Error

	if err != nil {
		return false, err
	}
	if u.ID == 0 {
		return false, nil
	}
	return true, nil
}

//检查是否已被注册。
//返回：是否查询到目标数据，以及任何错误。
func CheckRegisterOrNot(from string, u *ut.UserInformation) (bool, error) {
	err := dao.DB.Select("id").Where(u, from).Find(u).Error
	if err != nil {
		return false, err
	}
	if u.ID == 0 {
		return false, nil
	}
	return true, nil
}

//注册。
func PostRegisterInformation(u *ut.UserInformation) error {
	err := dao.DB.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}
