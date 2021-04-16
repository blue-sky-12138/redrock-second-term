package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	Issuer = "RedRock"
)

const (
	refreshToken = "refresh_token"
)

type JWTExt interface {
	apply(jwt *JWT)
}

type temJWTFunc func(j *JWT)

type funcJWT struct {
	f temJWTFunc
}

func (f funcJWT) apply(jwt *JWT) {
	f.f(jwt)
}

func WithAlg(alg string) funcJWT {
	return funcJWT{
		f: func(j *JWT) {
			j.Alg = alg
		},
	}
}

func WithIss(iss string) funcJWT {
	return funcJWT{
		f: func(j *JWT) {
			j.Issuer = iss
		}}
}

type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type UIfo struct {
	ID   uint   `json:"id"`
	Name string `json:"username"`
}

type Payload struct {
	//ID        int64       `json:"jti"`
	//Audience  string      `json:"aud"`
	Issuer   string `json:"iss"`
	IssuedAt string `json:"iat"`
	//NotBefore int64       `json:"nbf"`
	ExpiresAt string `json:"exp"`
	Subject   UIfo   `json:"sub"`
}

type JWT struct {
	Header
	Payload
	Signature string
	Token     string `json:"token"`
}

var (
	DefaultJWT = JWT{
		Header: Header{
			Typ: "JWT",
			Alg: "HS256",
		},
		Payload: Payload{
			Issuer: Issuer,
		},
	}
	HS256Key = []byte("RedRock12138")
)

func NewJWT(id uint, name string, je ...JWTExt) JWT {
	jwt := DefaultJWT
	for _, v := range je {
		v.apply(&jwt)
	}

	jwt.Payload.Subject.ID = id
	jwt.Payload.Subject.Name = name
	jwt.IssuedAt = strconv.FormatInt(time.Now().Unix(), 10)
	jwt.ExpiresAt = strconv.FormatInt(time.Now().Add(168*time.Hour).Unix(), 10)

	if jwt.Alg == "HS256" {
		jwt.HS256Encode()
	}

	return jwt
}

func DecodeJWT(token string) (JWT, error) {
	var j JWT

	lice := strings.Split(token, ".")
	if len(lice) != 3 {
		return JWT{}, fmt.Errorf("token Error")
	}

	dHead, err := base64.RawURLEncoding.DecodeString(lice[0])
	if err != nil {
		return JWT{}, err
	}

	err = json.Unmarshal(dHead, &j.Header)
	if err != nil {
		return JWT{}, err
	} else if j.Header.Typ != "JWT" || j.Header.Alg != "HS256" { //判断头部数据是否合法
		return JWT{}, fmt.Errorf("illegal type")
	}

	dPayload, err := base64.RawURLEncoding.DecodeString(lice[1])
	if err != nil {
		return JWT{}, err
	}

	err = json.Unmarshal(dPayload, &j.Payload)
	if err != nil {
		return JWT{}, err
	} else if j.Payload.Issuer != Issuer {
		return JWT{}, fmt.Errorf("illegal iss")
	} else if j.Payload.ExpiresAt < strconv.FormatInt(time.Now().Unix(), 10) {
		return JWT{}, fmt.Errorf("token is out of date")
	} else if j.Payload.IssuedAt >= j.Payload.ExpiresAt {
		return JWT{}, fmt.Errorf("illegal iat")
	} else if reflect.ValueOf(j.Payload.Subject).IsZero() {
		return JWT{}, fmt.Errorf("illegal sub")
	}

	j.HS256Encode()
	return j, nil
}

func (j *JWT) HS256Encode() {
	jHead, _ := json.Marshal(j.Header)
	jPayload, _ := json.Marshal(j.Payload)
	eHead := base64.RawURLEncoding.EncodeToString(jHead)
	ePayload := base64.RawURLEncoding.EncodeToString(jPayload)

	hash := hmac.New(sha256.New, HS256Key)
	hash.Write([]byte(eHead + "." + ePayload))
	sign := hash.Sum(nil)

	j.Signature = base64.RawURLEncoding.EncodeToString(sign)

	j.Token = eHead + "." + ePayload + "." + j.Signature
}

func TokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if len(h) < 7 {
			c.JSON(200, "token authentication failed")
			c.Abort()
		}

		token := h[7:]
		jwt, err := DecodeJWT(token)
		if err == nil && jwt.Token == token {
			c.Next()
		} else {
			c.JSON(200, "token authentication failed")
			c.Abort()
		}
	}
}

//相比于TokenCheck，这个函数不是中间件，而且会验证token里的name是否为refre_token，防止用token来更新token
func TokenRefreshChek(token string) bool {
	jwt, err := DecodeJWT(token)
	if err == nil && jwt.Token == token && jwt.Subject.Name == refreshToken {
		return true
	} else {
		return false
	}
}
