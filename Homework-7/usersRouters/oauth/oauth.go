package oauth

import (
	MySQL2 "SecondTerm/Homework-7/usersRouters/oauth/MySQL"
	jwt2 "SecondTerm/Homework-7/usersRouters/oauth/jwt"
	Redigo "SecondTerm/Homework-7/usersRouters/oauth/redigo"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Authorize struct {
	ResponseType string `json:"response_type" form:"response_type"`
	ClientId     int    `json:"client_id" form:"client_id"`
	RedirectUrl  string `json:"redirect_url" form:"redirect_url"`
	Scope        string `json:"scope" form:"scope"`
	State        string `json:"state" form:"state"`
}

type AuthCode struct {
	Code  string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}

type AccessToken struct {
	GrantType    string `json:"grant_type" form:"grant_type"`
	Code         string `json:"code" form:"code"`
	ClientId     int    `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	RedirectUri  string `json:"redirect_uri" form:"redirect_uri"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

const (
	stateKey          = "BlueSky"
	responseType      = "code"
	expiresIn         = 604800
	authorizationCode = "authorization_code"
	accessToken       = "access_token"
	refreshToken      = "refresh_token"
	tokenType         = "bearer"
	scope             = "read"
)

const (
	authorizeUrl = "http://localhost:8000/serve/oauth/authorize"
	callBackUrl  = "http://localhost:8000/serve/oauth/callback"
	tokenUrl     = "http://localhost:8000/serve/oauth/token"
)

var (
	pc *Redigo.Client
)

func NewAuthorize(clientId int, redirectUri string) Authorize {
	return Authorize{
		ResponseType: responseType,
		ClientId:     clientId,
		RedirectUrl:  redirectUri,
		Scope:        scope,
		State:        NewState(),
	}
}

func OAuthInit() error {
	//设置随机种子
	rand.Seed(time.Now().Unix())

	//redis初始化
	err := Redigo.RedigoInit()
	if err != nil {
		return fmt.Errorf("RedigoInit", err)
	}
	pc = Redigo.NewClient()

	//MySQL初始化
	err = MySQL2.MySQLInit()
	if err != nil {
		return fmt.Errorf("MySQLInit", err)
	}

	return nil
}

//OAuth验证请求
func OAuthRequest(clientId int, redirectUri string) {
	au := NewAuthorize(clientId, redirectUri)

	var i interface{} = au
	query := fillUrlValue(&i)
	resp, err := http.Get(authorizeUrl + "?" + query)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
}

//OAuth验证第一步
func OAuthAuthorize(c *gin.Context) {
	var au Authorize
	c.ShouldBind(&au)
	if au.ResponseType != responseType {
		c.JSON(200, "ResponseType Error")
		return
	}

	//验证client是否合法
	ci := MySQL2.ClientInfo{
		ClientID: au.ClientId,
	}
	if !ci.CheckClientID() {
		c.JSON(200, "client_id illegal")
		return
	}

	//获取code
	code := rand.Int()

	//将授权码存入redis
	pc.SetEX(600, code, strconv.Itoa(au.ClientId)+"."+au.RedirectUrl)

	//返回数据
	if au.RedirectUrl != "" {
		v := url.Values{}
		v.Add("code", strconv.Itoa(code))
		if au.State != "" {
			v.Add("state", au.State)
		}

		c.Redirect(http.StatusMovedPermanently, au.RedirectUrl+"?"+v.Encode())
	} else {
		var h gin.H
		if au.State != "" {
			h = gin.H{
				"code":  code,
				"state": au.State,
			}
		} else {
			h = gin.H{
				"code": code,
			}
		}

		c.JSON(200, h)
	}

}

//OAuth验证第二步
func OAuthCallBack(c *gin.Context) {
	var ac AuthCode
	c.ShouldBind(&ac)

	//如果数据不正确，直接丢弃
	if ac.Code == "" || ac.State == "" {
		return
	} else if !CheckState(ac.State) {
		return
	}

	//编写param
	at := AccessToken{
		GrantType:    authorizationCode,
		Code:         ac.Code,
		ClientId:     123456,
		ClientSecret: "RedRock",
		RedirectUri:  callBackUrl,
	}
	var i interface{} = at
	query := fillUrlValue(&i)

	//发送请求
	resp, err := http.Get(tokenUrl + "?" + query)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	//数据绑定
	var token Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		log.Println(err)
		return
	}

	//获取token后，设置cookie
	gh := gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	jsgh, _ := json.Marshal(gh)

	cookie := &http.Cookie{
		Name:     "user",
		Value:    base64.RawURLEncoding.EncodeToString(jsgh),
		MaxAge:   100000,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)
}

//OAuth验证第三步
func OAuthToken(c *gin.Context) {
	var at AccessToken
	c.ShouldBind(&at)

	//验证信息是否缺失
	if at.Code == "" || at.ClientId == 0 || at.ClientSecret == "" {
		c.JSON(200, "data empty")
		return
	}

	//验证信息是否合法
	if at.GrantType == authorizationCode {
		get, _ := pc.Get(at.Code)
		if get == nil {
			c.JSON(200, "code is error or is out of date")
			return
		} else {
			s, _ := redis.String(get, nil)
			pc.Del(at.Code)

			lice := strings.Split(s, ".")
			if len(lice) == 2 {
				if lice[0] != strconv.Itoa(at.ClientId) || lice[1] != at.RedirectUri {
					c.JSON(200, "data illegal")
					return
				}
			} else {
				if lice[0] != strconv.Itoa(at.ClientId) || "" != at.RedirectUri {
					c.JSON(200, "data illegal")
					return
				}
			}
		}
	} else if at.GrantType == refreshToken {
		if !jwt2.TokenRefreshChek(at.Code[7:]) {
			c.JSON(200, "refresh token error")
			return
		}
	} else {
		c.JSON(200, "data grant_type error")
		return
	}

	//验证身份
	ci := MySQL2.ClientInfo{
		ClientID:     at.ClientId,
		ClientSecret: at.ClientSecret,
	}
	if !ci.CheckClientInfo() {
		c.JSON(200, "illegal identity")
		return
	}

	//处理请求
	atk := Token{
		AccessToken:  jwt2.NewJWT(uint(at.ClientId), accessToken).Token,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		RefreshToken: jwt2.NewJWT(uint(at.ClientId), refreshToken).Token,
		Scope:        scope,
	}
	c.JSON(200, atk)
}

//新建一个state
//以时间种子的随机数作为密钥和数据
func NewState() string {
	key := []byte(strconv.Itoa(rand.Int()))

	h := hmac.New(md5.New, []byte(stateKey))
	h.Write(key)
	sum := h.Sum(nil)

	eKey := base64.RawURLEncoding.EncodeToString(key)
	eSum := base64.RawURLEncoding.EncodeToString(sum)

	if pc.Exists(eKey) == 1 {
		return NewState()
	}

	pc.SetEX(300, eKey, eSum)

	return eKey + "." + eSum
}

//检查state是否正确
func CheckState(state string) bool {
	lice := strings.Split(state, ".")
	if len(lice) != 2 {
		return false
	}

	get, err := pc.Get(lice[0])
	if err != nil {
		return false
	}
	pc.Del(lice[0])
	res, _ := redis.String(get, nil)

	return lice[1] == res
}

//自动填充form参数
func fillUrlValue(data *interface{}) string {
	typ := reflect.TypeOf(*data)
	val := reflect.ValueOf(*data)

	res := ""
	fieldNum := typ.NumField()

	for i := 0; i < fieldNum; i++ {
		t := typ.Field(i)
		v := val.Field(i)

		//以form标签作为key来加入到values中
		if v.Kind() == reflect.Int {
			res += t.Tag.Get("form") + "=" + strconv.FormatInt(v.Int(), 10) + "&"
		} else {
			res += t.Tag.Get("form") + "=" + v.String() + "&"
		}
	}

	return res[:len(res)-1]
}

//类似于fillUrlValue，自动填充form参数，但会忽略零值
func fillUrlValueWithoutZero(data *interface{}) string {
	typ := reflect.TypeOf(*data)
	val := reflect.ValueOf(*data)

	res := ""
	fieldNum := typ.NumField()

	for i := 0; i < fieldNum; i++ {
		t := typ.Field(i)
		v := val.Field(i)

		if v.IsZero() {
			continue
		}

		//以form标签作为key来加入到values中
		if v.Kind() == reflect.Int {
			res += t.Tag.Get("form") + "=" + strconv.FormatInt(v.Int(), 10) + "&"
		} else {
			res += t.Tag.Get("form") + "=" + v.String() + "&"
		}
	}

	return res[:len(res)-1]
}
