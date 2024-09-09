package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	// "golang.org/x/net/websocket"
)

// 获取所有用户
// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功，-1失败
		"message": "查找成功！",
		"data":    data,
	})
}

//false代表参数可选

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	fmt.Println(user.Name, " >>>>>>>", password, repassword)
	salt := fmt.Sprintf("%06d", rand.Int31())
	data := models.FindUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名或密码不能为空",
			"data":    user,
		})
		return
	}
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "用户名已经注册!",
			"data":    user,
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}
	user.PassWord = utils.MakePassWord(password, salt)
	user.Salt = salt
	fmt.Println(user.PassWord)
	user.LoginTime = time.Now()
	user.LoginOutTime = time.Now()
	user.HeartbeatTime = time.Now()
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功，-1失败
		"message": "新增用户成功",
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功，-1失败
		"message": "删除用户成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Avatar = c.PostForm("icon")
	user.Email = c.PostForm("email")

	fmt.Println("update :", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "修改参数不匹配",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0, //0成功，-1失败
			"message": "修改用户成功",
			"data":    user,
		})
	}

}

// 登录验证
// FindUserByNameAndPwd
// @Summary 登录验证
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/FindUserByNameAndPwd [Post]
func FindUserByNameAndPwd(c *gin.Context) {
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	fmt.Println(name, password)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "该用户不存在",
			"data":    user,
		})
		return
	}

	flag := utils.VaildPassWord(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //0成功，-1失败
			"message": "密码不正确",
			"data":    user,
		})
		return
	}
	pwd := utils.MakePassWord(password, user.Salt)
	//登录验证
	data := models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, //0成功，-1失败
		"message": "登录成功",
		"data":    data,
	})
}

// 定义了一个 WebSocket 的升级器 upGrade，
// 用于升级 HTTP 请求为 WebSocket 连接。
// 防止跨域站点的伪造请求
// var upGrade = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// Gin 的请求处理函数，用于处理 WebSocket 连接的升级和消息发送。
// func SendMsg(c *gin.Context) {
// 	// 将 HTTP 连接升级为 WebSocket 连接。
// 	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
// 	// 如果升级失败，会打印错误并返回。
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	//延迟关闭连接
// 	defer func(ws *websocket.Conn) {
// 		err = ws.Close()
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}(ws)
// 	//调用 MsgHandler 函数处理 WebSocket 的消息发送
// 	MsgHandler(ws, c)
// }

// // 用于处理 WebSocket 的消息发送。
// // 使用一个无限循环来不断地从 Redis 订阅消息
// // 并将消息发送到 WebSocket 连接
//
//	func MsgHandler(ws *websocket.Conn, c *gin.Context) {
//		for {
//			msg, err := utils.Subscribe(c, utils.PublishKey)
//			if err != nil {
//				fmt.Println(err)
//			}
//			tm := time.Now().Format("2006-01-02 15:04:06")
//			m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
//			err = ws.WriteMessage(1, []byte(m))
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//	}
//
// postform用于POST请求中的表单数据
// formvalue用于GET与POST请求中的表单数据 POST优先级高
// 从redis中获取两者之间的聊天记录
func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

// 用户准备聊天
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// 将用户的好友列表返回给前端
func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))
	utils.RespOKList(c.Writer, users, len(users))
}

// 用户根据名字添加好友
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 用户新建一个群
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加载用户的群列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加入群聊 通过群id加群
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	data, msg := models.JoinGroups(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 根据id查找用户
func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}
