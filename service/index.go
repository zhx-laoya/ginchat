// 跳转页面的服务+聊天建立连接服务
package service

import (
	"ginchat/models"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "welcome!!",
	// })
}

// post请求
func ToRegister(c *gin.Context) {
	///解析文件
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	//register为传递給模板的函数
	ind.Execute(c.Writer, "register")
}

// 用户发送进入聊天页面的请求，将这些模板文件返回给前端
func ToChat(c *gin.Context) {
	//解析文件
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}

	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	//把id和token传递给前端
	ind.Execute(c.Writer, user)
}

// 用户发送信息，请求建立连接进行聊天
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
