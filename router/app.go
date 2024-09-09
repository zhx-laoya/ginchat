package router

import (
	// "ginchat/docs"
	"ginchat/service"

	"github.com/gin-gonic/gin"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	// docs.SwaggerInfo.BasePath = ""
	// //swagger用于生成接口
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//加载静态资源
	//静态资源
	r.Static("/asset", "asset/")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	r.LoadHTMLGlob("views/**/*")

	//页面跳转
	r.GET("/", service.GetIndex) //首页
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)        //注册界面
	r.GET("/toChat", service.ToChat)                //聊天页面
	r.GET("/chat", service.Chat)                    //聊天
	r.POST("/searchFriends", service.SearchFriends) //朋友列表

	//用户使用模块
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.POST("/user/find", service.FindByID) //根据id查找用户

	//发送消息
	// r.GET("/user/sendMsg", service.SendMsg)
	//发送消息
	// r.GET("/user/sendUserMsg", service.SendUserMsg)
	//添加好友
	r.POST("/contact/addfriend", service.AddFriend)
	//上传文件
	r.POST("/attach/upload", service.Upload)
	//创建群
	r.POST("/contact/createCommunity", service.CreateCommunity)
	//获取群列表
	r.POST("/contact/loadcommunity", service.LoadCommunity)
	//加群
	r.POST("/contact/joinGroup", service.JoinGroups)
	//心跳续命 不合适  因为Node  所以前端发过来的消息再receProc里面处理
	// r.POST("/user/heartbeat", service.Heartbeat)
	//将信息存储到redis中
	r.POST("/user/redisMsg", service.RedisMsg)
	return r
}
