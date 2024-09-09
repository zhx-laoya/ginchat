// md5加密
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempstr := h.Sum(nil)
	return hex.EncodeToString(tempstr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 随机数加密
func MakePassWord(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func VaildPassWord(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
