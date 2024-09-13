package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func SplicingString(str []string, s string) (newstr string) {
	newstr = strings.Join(str, s)
	return
}

func DingSecret(secret string) (sign string) {
	// 1. 获取当前时间戳，单位是毫秒
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	// 2. 把时间戳和密钥拼接成字符串
	stringToSign := timestamp + "\n" + secret

	// 3. 使用 HMAC-SHA256 进行加密
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	signData := h.Sum(nil)

	// 4. 进行 Base64 编码
	signature := base64.StdEncoding.EncodeToString(signData)

	// 5. 进行 URL 编码
	signature = url.QueryEscape(signature)
	sign = "&timestamp=" + timestamp + "&sign=" + signature
	return
}
