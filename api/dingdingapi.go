package api

import (
	init_tool "bug-notify/init-tool"
	"bug-notify/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendMessage(data model.SendMsg) error {
	URL := init_tool.Conf.Robot.URL
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

	msg["msgtype"] = data.MsgType
	if data.MsgType == "markdown" {
		msg[data.MsgType] = map[string]interface{}{
			"title": "bug",
			"text":  data.Content,
		}
	} else {
		msg[data.MsgType] = map[string]interface{}{
			"title":       "bug",
			"text":        data.Content,
			"singleTitle": "问题详情",
			"singleURL":   data.Url,
		}
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
