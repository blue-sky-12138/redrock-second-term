package model

import (
	"SecondTerm/Homework-7/usersCenter/dao"
	ut "SecondTerm/Homework-7/usersCenter/utilities"
	"gorm.io/gorm"
	"regexp"
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

	err = dao.DB.Select("id", "name", "password", "md5salt").Where(u, from).Find(u).Error

	if err != nil {
		ut.LogError("GetUserPassword Error", err)
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
		ut.LogError("CheckRegisterOrNot Error", err)
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

////更新生日。
//func ChangeUserBirthday(userId uint, detail int64) error{
//	u := &ut.UserInformation{
//		Model: gorm.Model{ID: userId},
//		Birthday: ,
//	}
//	err := dao.ChangeUserInformation("birthday", u)
//	if err != nil {
//	return err
//	}
//	return nil
//}

//更新签名。
func ChangeUserSignature(userId uint, detail string) error {
	u := &ut.UserInformation{
		Model:    gorm.Model{ID: userId},
		HeadPath: detail,
	}
	err := dao.ChangeUserInformation("signature", u)
	if err != nil {
		return err
	}
	return nil
}

//更新性别。
func ChangeUserGender(userId uint, detail int) error {
	u := &ut.UserInformation{
		Model:  gorm.Model{ID: userId},
		Gender: detail,
	}
	err := dao.ChangeUserInformation("gender", u)
	if err != nil {
		return err
	}
	return nil
}

//更新昵称。
func ChangeUserNickname(userId uint, detail string) error {
	u := &ut.UserInformation{
		Model: gorm.Model{ID: userId},
		Name:  detail,
	}
	err := dao.ChangeUserInformation("nickname", u)
	if err != nil {
		return err
	}
	return nil
}

//更新邮箱。
func ChangeUserEmail(userId uint, detail string) error {
	u := &ut.UserInformation{
		Model: gorm.Model{ID: userId},
		Email: detail,
	}
	err := dao.ChangeUserInformation("email", u)
	if err != nil {
		return err
	}
	return nil
}

//更新手机号。
func ChangeUserTelephone(userId uint, detail int64) error {
	u := &ut.UserInformation{
		Model:     gorm.Model{ID: userId},
		Telephone: detail,
	}
	err := dao.ChangeUserInformation("telephone", u)
	if err != nil {
		return err
	}
	return nil
}

//更新头像
func ChangeUserHead(userId uint, detail string) error {
	u := &ut.UserInformation{
		Model:    gorm.Model{ID: userId},
		HeadPath: detail,
	}
	err := dao.ChangeUserInformation("head_path", u)
	if err != nil {
		return err
	}
	return nil
}

//获取用户头像路径
func GetUserHeadPath(userId uint) (string, error) {
	u := &ut.UserInformation{
		Model: gorm.Model{ID: userId},
	}
	err := dao.DB.Select("head_path").Find(u).Error
	if err != nil {
		ut.LogError("GetUserHeadPath Error", err)
		return "", err
	}
	return u.HeadPath, nil
}
