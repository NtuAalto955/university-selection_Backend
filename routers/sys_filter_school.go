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
	errCode := 200

	defer func() {
		c.Set("errorCode", errCode)
	}()

	req := &sysRequest.SchoolFilterReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		errCode = 400
		c.JSON(errCode, gin.H{
			"success": false,
			"msg":     "输入参数错误",
			"result":  err,
		})
		return
	}
	res, err := biz.RefSchoolFilterBiz().FilterSchool(req)
	if err != nil {
		errCode = 400
		c.JSON(errCode, gin.H{
			"success": false,
			"msg":     "学校筛选失败",
			"result":  err,
		})
		return
	}

	resJson, err := json.Marshal(res)

	if err != nil {
		global.GLog.Error(fmt.Sprintf("Marshal error:%v", err))
		errCode = 400
		c.JSON(errCode, gin.H{
			"success": false,
			"msg":     "学校筛选失败",
			"result":  err,
		})
		return
	}
	errCode = 200
	c.JSON(errCode, gin.H{
		"success": true,
		"msg":     "filter success",
		"result":  string(resJson),
	})
}
