package main

import (
	"admin_project/biz"
	"admin_project/core"
	_ "admin_project/docs"
	"admin_project/global"
	"admin_project/middlerware"
	"admin_project/routers"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @termsOfService http://swagger.io/terms/
func main() {
	//启动日志
	global.GLog = core.Zap()
	global.GLog.Debug("server running")
	//启动配置读取
	global.G_Viper = core.Viper()
	//连接数据库
	global.G_DB = core.Db()
	//global.G_DB.AutoMigrate(&global.User{}, &global.Comment{}, &global.OfferInfo{})
	db, _ := global.G_DB.DB()

	defer db.Close()
	//u := User{Password: "test",Username: "test4"}
	//gDb.Create(&u)

	//
	//admin 服务器启动
	https := gin.Default()
	http := gin.Default()
	//启动接口文档swagger
	//https.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//fmt.Println("在线api文档部署在：http://localhost:8080/swagger/index.html")
	// 微信获取access token
	biz.GetAccessToken()
	//公共路由 注册，登录，验证码
	publicRouter := https.Group("")
	publicRouter.Use(middlerware.TlsHandler())
	publicRouter.Use(middlerware.InjectCtx)

	{
		publicRouter.GET("/captcha", routers.Captcha, middlerware.ReportProm)
		publicRouter.POST("/register", routers.RegisterHandler, middlerware.ReportProm)
		publicRouter.POST("/login", routers.LoginHandler, middlerware.ReportProm)
		publicRouter.POST("/get_filter_school", routers.FilterSchoolHandler(), middlerware.ReportProm)
		publicRouter.GET("/send_wx_msg", routers.VerifyWxToken(), middlerware.ReportProm)
		publicRouter.POST("/send_wx_msg", routers.ProcessWxMsgHandler(), middlerware.ReportProm)

	}

	httpRouter := http.Group("/")
	{
		// 上报promethus
		httpRouter.GET("metrics", gin.WrapH(promhttp.Handler()), middlerware.ReportProm)
	}
	//用户路由   访问前需要认证token
	//usrRouter := s.Group("")
	//
	//usrRouter.Use(middlerware.Auth)
	//{
	//	usrRouter.GET("userinfo", routers.GetinfoHandler)
	//	usrRouter.POST("deleteUser", routers.DeleteUserHandler)
	//	usrRouter.POST("changepassword", routers.ChangePassword)
	//
	//	s.POST("/addcomment", routers.AddComment)
	//	s.POST("/deletecomment", routers.DeleteComment)
	//	s.GET("/getcomment", routers.GetComment)
	//}

	// 服务启动在443端口,使用https，
	go func() {
		if err := https.RunTLS(":443", "ssl.pem", "ssl.key"); err != nil {
			global.GLog.Error("https server is fail!")
		}
	}()

	if err := http.Run(":8080"); err != nil {
		global.GLog.Error("http server is fail!")

	}
}

// ShowAccount godoc
// @Summary Show an account
// @Tags Example API
// @Description get string by ID
// @Produce  json
// @Success 200
// @Router /ping [get]
func pang(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
