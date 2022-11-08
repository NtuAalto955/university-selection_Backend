package middlerware

import (
	"admin_project/global"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"time"
)

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":443",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}

func InjectCtx(c *gin.Context) {
	c.Set("processingTime", time.Now())
	c.Set("errorCode", 0)
}

func Auth(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(405, gin.H{
			"success": false,
			"msg":     "you should login first!",
		})
		c.Abort()
		return
	}
	username, err := PhaseToken(tokenString)
	if username == "" || err != nil {
		c.JSON(405, gin.H{
			"success": false,
			"msg":     "auth failed!",
		})
		c.Abort()
		return
	}
	global.GLog.Info(username + " has visited!")

}
