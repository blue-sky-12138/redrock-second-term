package utilities

import (
	"gorm.io/gorm"
)

func (u UserInformation) TableName() string {
	return "users_information"
}

//用户信息
type UserInformation struct {
	gorm.Model
	Name      string `json:"username" form:"username" xml:"username"` //用户名(同时这也是接受前端传入的登录信息的变量)
	Password  string `json:"password" form:"password" xml:"password"`
	Md5salt   string `json:"md5_salt"`
	Telephone int64  `json:"telephone" form:"telephone" xml:"telephone"`
	Email     string `json:"email" form:"email" xml:"email"`
	HeadPath  string `json:"head_path"`
	Signature string `json:"signature"`
	Gender    int    `json:"gender"`
	Birthday  string `json:"birthday"`
	Vip       int    `json:"vip"`
	Level     int    `json:"level"`
}
