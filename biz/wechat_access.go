package biz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"strconv"
	"time"
)

var (
	AppID       = "wx98b4bdded43c3177"
	AppSecrect  = "01958b0f88d702375e84c6b2ca1054d0"
	AccessToken = ""
)

type recvMsg struct {
	AccessCode string `json:"access_code"`
	ExpireTime string `json:"expires_in"`
}

func loop() int64 {
	session, err := http.NewRequest("GET", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", AppID, AppSecrect), nil)
	resp, err := http.DefaultClient.Do(session)
	if err != nil {
		log.Fatalln(err)
		return 0
	}
	defer resp.Body.Close()
	var accessToken recvMsg
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&accessToken)
	if err != nil {
		log.Fatalln(err)
		return 0
	}
	expirTime, _ := strconv.ParseInt(accessToken.ExpireTime, 10, 64)
	AccessToken = accessToken.AccessCode
	return expirTime
}
func GetAccessToken() {
	expirTime := loop()
	if expirTime <= 0 {
		return
	}
	t := time.NewTimer(time.Duration(expirTime-60) * time.Second)
	for {
		select {
		case <-t.C:
			expirTime = loop()
			t.Reset(time.Duration(expirTime-60) * time.Second)
		}
	}
}

type ReplyMsg struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// 微信公众号回复有5s限制，需要转换第三方客服异步回复
func ThirdPartyReply(userOpenID string, msg string) {
	reply := &ReplyMsg{
		Touser:  userOpenID,
		Msgtype: "text",
		Text: struct {
			Content string `json:"content"`
		}{msg},
	}
	str, err := json.Marshal(reply)
	if err != nil {
		log.Fatalln(err)
		return
	}
	session, err := http.NewRequest("POST", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%v", AccessToken), bytes.NewReader(str))
	resp, err := http.DefaultClient.Do(session)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer resp.Body.Close()
}
