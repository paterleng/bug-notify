package api

import (
	"bug-notify/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
			"atMobiles": data.AtMobiles,
			"isAtAll":   data.IsAtAll,
		}
	}
	msg["msgtype"] = "markdown"
	msg["markdown"] = map[string]interface{}{
		"title": "bug",
		"text":  data.Content,
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
