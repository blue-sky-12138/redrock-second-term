package dao

import (
	ut "SecondTerm/Homework-7/fileRouters/utilities"
)

//更新用户信息。
func ChangeUserInformation(focus string, u *ut.UserInformation) error {
	exist, err := checkIfExist(u)
	if err != nil {
		ut.LogError("checkIfExist Error", err)
		return err
	} else if !exist {
		return ut.ErrorUserNotExist
	}

	err = DB.Model(u).Where("id = ?", u.ID).Updates(u).Error
	if err != nil {
		ut.LogError("ChangeUserInformation Error", err)
		return err
	}
	return nil
}

func checkIfExist(u *ut.UserInformation) (bool, error) {
	temId := u.ID
	u.ID = 0
	err := DB.Select("id").Where("id = ?", temId).Find(u).Error
	if err != nil {
		return false, err
	}
	if u.ID == 0 {
		return false, nil
	}
	return true, nil
}
