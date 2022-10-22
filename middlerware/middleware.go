package middlerware

import (
	"admin_project/global"
	"github.com/gin-gonic/gin"
	"time"
)

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
