package api

import (
	"bug-notify/model"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetUserIDByPhone(phone string) (userid string, err error) {
	var client *http.Client
	var request *http.Request
	var resp *http.Response
	var body []byte
	URL := "https://oapi.dingtalk.com/topapi/v2/user/getbymobile?access_token=" + ""
	client = &http.Client{Transport: &http.Transport{ //对客户端进行一些配置
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}, Timeout: time.Duration(time.Second * 5)}
	//此处是post请求的请求题，我们先初始化一个对象
	b := struct {
		Mobile string
	}{
		Mobile: phone,
	}
	//然后把结构体对象序列化一下
	bodymarshal, err := json.Marshal(&b)
	if err != nil {
		return
	}
	//再处理一下
	reqBody := strings.NewReader(string(bodymarshal))
	//然后就可以放入具体的request中的
	request, err = http.NewRequest(http.MethodPost, URL, reqBody)
	if err != nil {
		return
	}
	resp, err = client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body) //把请求到的body转化成byte[]
	if err != nil {
		return
	}
	r := struct {
		model.DingResponseCommon
		Result struct {
			UserID string `json:"userid"`
		} `json:"result"`
	}{}
	//把请求到的结构反序列化到专门接受返回值的对象上面
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		return "", errors.New(r.Errmsg)
	}
	// 此处举行具体的逻辑判断，然后返回即可

	return r.Result.UserID, nil
}

func SendMessage(data model.SendMsg) error {
	URL := "https://oapi.dingtalk.com/robot/send?access_token=8ff6cde9a01910e897cb6461e75bd515ed9d683cc4924aad9439fda3e9689de1"
	b := []byte{}
	msg := map[string]interface{}{}
	//@的人
	if data.IsAtAll {
		msg["at"] = map[string]interface{}{
			"isAtAll": data.IsAtAll,
		}
	} else {
		msg["at"] = map[string]interface{}{
			"atUserIds": []string{data.AtUserID},
			"isAtAll":   data.IsAtAll,
		}
	}

	msg = map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": data.Content,
		},
	}
	b, _ = json.Marshal(msg)
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	date, err := ioutil.ReadAll(resp.Body)
	r := model.DingResponseCommon{}
	err = json.Unmarshal(date, &r)
	if err != nil {
		return err
	}
	if r.Errcode != 0 {
		fmt.Println(r.Errmsg)
		return errors.New(r.Errmsg)
	}
	return nil
}
