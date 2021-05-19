package utilities

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strconv"
	"time"
)

//MD5即时加密的快捷方式
func CryptographyNow(Data string) (string, string) {
	Md5salt := strconv.FormatInt(time.Now().Unix(), 10)
	return Cryptography(Data, Md5salt), Md5salt
}

//MD5加密
func Cryptography(Data string, Md5salt string) string {
	has := md5.New()
	io.WriteString(has, Data+Md5salt)
	tem := has.Sum(nil)
	Result := hex.EncodeToString(tem)
	return Result
}
