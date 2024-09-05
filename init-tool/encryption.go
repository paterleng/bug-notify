package init_tool

import (
	"golang.org/x/crypto/bcrypt"
)

func Encryption(str string) (string, error) {
	//加密密码
	b, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Contrast(old, new string) (bool, error) {
	//对比密码
	err := bcrypt.CompareHashAndPassword([]byte(old), []byte(new))
	if err != nil {
		return false, err
	}
	return true, nil
}
