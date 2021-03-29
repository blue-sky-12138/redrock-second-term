package utilities

import (
	"regexp"
	"strconv"
)

//手机号正确性检查.
func TelephoneCheck(number int64) bool {
	reg := regexp.MustCompile("(13|14|15|17|18|19)[0-9]{9}")
	if reg.MatchString(strconv.FormatInt(number, 10)) {
		return true
	} else {
		return false
	}
}

//用户名有效性检查.
//只能包括数字、字母、汉字的组合，且不能只包含数字.
func UserNameCheck(name string) bool {
	reg := regexp.MustCompile("[^A-Za-z0-9\u4e00-\u9fa5]")
	if !reg.MatchString(name) {
		reg = regexp.MustCompile("[^0-9]")
		if reg.MatchString(name) {
			return true
		}
		return false
	} else {
		return false
	}
}

//邮箱正确性检测.
func EmailCheck(email string) bool {
	reg := regexp.MustCompile("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}")
	if reg.MatchString(email) {
		return true
	} else {
		return false
	}
}

//密码规范性检查.
//最少8位,最多16位，只能包括数字、字母、下划线.
func PasswordCheck(password string) bool {
	reg := regexp.MustCompile("[^\\w]")
	if len(password) >= 8 && len(password) <= 16 && !reg.MatchString(password) {
		return true
	} else {
		return false
	}
}
