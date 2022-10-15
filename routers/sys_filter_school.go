package routers

import (
	"admin_project/biz"
	"admin_project/global"
	"admin_project/sysRequest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func FilterSchoolHandler(c *gin.Context) {
	req := &sysRequest.SchoolFilterReq{}

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"msg":     "输入参数错误",
			"result":  err,
		})
		return
	}
	res, err := biz.RefSchoolFilterBiz().FilterSchool(req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"msg":     "学校筛选失败",
			"result":  err,
		})
		return
	}

	resJson, err := json.Marshal(res)

	if err != nil {
		global.GLog.Error(fmt.Sprintf("Marshal error:%v", err))
		c.JSON(400, gin.H{
			"success": false,
			"msg":     "学校筛选失败",
			"result":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"msg":     "filter success",
		"result":  string(resJson),
	})
}
