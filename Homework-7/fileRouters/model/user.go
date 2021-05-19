package model

import (
	"SecondTerm/Homework-7/fileRouters/dao"
	ut "SecondTerm/Homework-7/fileRouters/utilities"
	"gorm.io/gorm"
)

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
