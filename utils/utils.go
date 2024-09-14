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
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	stringToSign := timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	signData := h.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(signData)
	signature = url.QueryEscape(signature)
	sign = "&timestamp=" + timestamp + "&sign=" + signature
	return
}
