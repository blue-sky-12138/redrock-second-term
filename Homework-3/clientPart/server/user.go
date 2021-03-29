package server

import (
	"bytes"
	ut "clientPart/utilities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type UserInfo struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

func Login() bool {
	var (
		u UserInfo
	)

	fmt.Printf("输入用户名：")
	fmt.Scanln(&u.Name)
	fmt.Printf("输入密码：")
	fmt.Scanln(&u.Password)

	err, resp := PostWithJson(LoginAddr, &u)
	if err != nil {
		return false
	}

	if resp.Code != 200 {
		fmt.Println(resp.Message)
		return false
	}

	UserName = u.Name
	return true
}

//用户注册
func Register() bool {
	var (
		u UserInfo
	)

	fmt.Printf("输入用户名(只能包括数字、字母、汉字的组合，且不能只包含数字)：")
	fmt.Scanln(&u.Name)
	fmt.Printf("输入密码(最少8位,最多16位，只能包括数字、字母、下划线)：")
	fmt.Scanln(&u.Password)

	err, resp := PostWithJson(RegAddr, &u)
	if err != nil {
		return false
	}

	if resp.Code != 300 {
		fmt.Println(resp.Message)
		return false
	}

	UserName = u.Name
	return true
}

func PostWithJson(addr string, v interface{}) (error, *ut.Resp) {
	var (
		res ut.Resp
	)

	body, err := json.Marshal(&v)
	header := bytes.NewReader(body)
	resp, err := http.Post(addr, "application/json", header)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return err, nil
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, &res
}
