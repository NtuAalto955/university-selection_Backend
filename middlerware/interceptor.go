package middlerware

import (
	"admin_project/monitor"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func ReportProm(c *gin.Context) {
	errCode, exist := c.Get("errorCode")
	if !exist {
		errCode = 0
	}
	startTime, exist := c.Get("processingTime")
	if !exist {
		startTime = time.Now()
	}
	costTime := time.Since(startTime.(time.Time)).Seconds()
	monitor.RefPromMonitor().ReportHttpCounter(c.Request.RequestURI, strconv.Itoa(errCode.(int)))
	monitor.RefPromMonitor().ReportHttpHistogram(c.Request.RequestURI, strconv.Itoa(errCode.(int)), costTime)

}
