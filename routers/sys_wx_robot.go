package routers

import (
	"admin_project/biz"
	"admin_project/util"
	"admin_project/util/wxbizmsgcrypt"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/groupcache/singleflight"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	encodingAESKey = "aFHzF817JMYATUSAQPjhC68vRWpyZ5bU7a6bQuqh98C"
	to_xml         = `<xml><ToUserName><![CDATA[oia2TjjewbmiOUlr6X-1crbLOvLw]]></ToUserName><FromUserName><![CDATA[gh_7f083739789a]]></FromUserName><CreateTime>1407743423</CreateTime><MsgType>  <![CDATA[video]]></MsgType><Video><MediaId><![CDATA[eYJ1MbwPRJtOvIEabaxHs7TX2D-HV71s79GUxqdUkjm6Gs2Ed1KF3ulAOA9H1xG0]]></MediaId><Title><![CDATA[testCallBackReplyVideo]]></Title><Descript  ion><![CDATA[testCallBackReplyVideo]]></Description></Video></xml>`
	token          = "anXdXcBqzVUPIMxsPJnTnuVlscdMIW"
	nonce          = "1320562132"
	appid          = "qLwDXki7PzATAi2"
	reqGroup       singleflight.Group
)

func ProcessWxMsgHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取用户信息
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		recv := util.ToTextMsg(bodyBytes)

		// 处理信息
		// 回复消息
		replyMsg := ":"

		// 关注公众号事件
		if recv.MsgType == "event" {
			if recv.Event == "unsubscribe" {
				biz.DefaultGPT.DeleteUser(recv.FromUserName)
			}
			if recv.Event != "subscribe" {
				util.TodoEvent(c.Writer)
				return
			}
			replyMsg = ":) 感谢你发现了这里"
		} else if recv.MsgType == "text" {
			// 【收到不支持的消息类型，暂无法显示】
			if strings.Contains(recv.Content, "【收到不支持的消息类型，暂无法显示】") {
				util.TodoEvent(c.Writer)
				return
			}
			// chatGPT 基本都会超时，调用第三方客服回复，48h只有20条
			go func() {
				ProcessGptMsg(recv)
			}()
		} else {
			util.TodoEvent(c.Writer)
			return
		}

		textRes := &util.TextRes{
			ToUserName:   recv.FromUserName,
			FromUserName: recv.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "transfer_customer_service",
			Content:      replyMsg,
		}
		_, err := c.Writer.Write(textRes.ToXml())
		if err != nil {
			log.Fatalln(err)
		}
	}
}
func ProcessGptMsg(recv *util.TextMsg) {
	_, err := reqGroup.Do(strconv.FormatInt(recv.MsgId, 10), func() (interface{}, error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		select {
		case msg := <-biz.DefaultGPT.SendMsgChan(recv.Content, recv.FromUserName, ctx):
			biz.ThirdPartyReply(recv.FromUserName, msg)
			return nil, nil
		case <-time.After(300*time.Second + 500*time.Millisecond):
			// 超时返回错误
			return "", fmt.Errorf("请求超时, MsgId: %d", recv.MsgId)
		}
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
}

//验证 微信token
func VerifyWxToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		errCode := 200

		defer func() {
			c.Set("errorCode", errCode)
		}()
		r := c.Request
		fmt.Println(r.RequestURI, r.Header, "1", c.Request.PostForm, c.PostForm("nonce"))
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		signature := c.Query("signature")
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)

		fmt.Println(string(bodyBytes))
		fmt.Println(timestamp, nonce, signature)
		cryptor, _ := wechataes.NewWechatCryptor(appid, token, encodingAESKey)
		msg, cryptErr := cryptor.DecryptMsg(signature, timestamp, nonce, string(bodyBytes))
		fmt.Println(msg, cryptErr)

	}
}
