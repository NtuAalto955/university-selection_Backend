package routers

import (
	"admin_project/util/wxbizmsgcrypt"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

var (
	encodingAESKey = "aFHzF817JMYATUSAQPjhC68vRWpyZ5bU7a6bQuqh98C"
	to_xml         = `<xml><ToUserName><![CDATA[oia2TjjewbmiOUlr6X-1crbLOvLw]]></ToUserName><FromUserName><![CDATA[gh_7f083739789a]]></FromUserName><CreateTime>1407743423</CreateTime><MsgType>  <![CDATA[video]]></MsgType><Video><MediaId><![CDATA[eYJ1MbwPRJtOvIEabaxHs7TX2D-HV71s79GUxqdUkjm6Gs2Ed1KF3ulAOA9H1xG0]]></MediaId><Title><![CDATA[testCallBackReplyVideo]]></Title><Descript  ion><![CDATA[testCallBackReplyVideo]]></Description></Video></xml>`
	token          = "anXdXcBqzVUPIMxsPJnTnuVlscdMIW"
	nonce          = "1320562132"
	appid          = "qLwDXki7PzATAi2"
)

func SendWxMsgHandler() gin.HandlerFunc {
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
