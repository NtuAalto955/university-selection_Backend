package routers

import (
	"admin_project/biz"
	"admin_project/global"
	"admin_project/sysRequest"
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/buffer"
)

func FilterSchoolHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		retJson, err := json.Marshal(res)
		if err != nil {
			global.GLog.Error(fmt.Sprintf("marshal error,err:%v", err))
		}
		w := &buffer.Buffer{}
		zipWriter := zip.NewWriter(w)
		file, _ := zipWriter.Create("result.txt")
		file.Write(retJson)
		zipWriter.Close()
		errCode = 200

		c.JSON(errCode, gin.H{
			"success": true,
			"msg":     "filter success",
			"result":  w.String(),
		})
	}
}
