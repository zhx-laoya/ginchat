//这一块实现了消息的传递，私聊群聊，心跳检测
// 注意：A和B的实时通信，本质上是A和B与服务端发websocket连接，再通过服务端处理后将信息传播给B，而不是直接连接

// 消息表 后面可以用redis来存放
// 发送消息
// 需要：发送者id，接收者id，消息类型，发送的内容，发送类型
// 接受消息：携程+管道
// 消息传递：先进入chat,用户绑定websocket连接-》进入接受与发送逻辑
// 接受逻辑分为 私聊，群聊，心跳
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserId     int64  // 发送者
	TargetId   int64  //接收者
	Type       int    //发送类型 1私聊 2群聊 3心跳
	Media      int    //消息类型 1文字 2表情包 3语音 4图片
	Content    string //消息内容
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string //图片
	Url        string //URL相关
	Desc       string //描述
	Amount     int    //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn          *websocket.Conn //websocket连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息管道
	GroupSets     set.Interface   //好友/群
}

// 需要重写此方法才能完整的msg转byte[] 注意byte[]转msg用json.unmarshal即可
func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// 映射关系,用户对应于某个连接
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 当客户端发送一个 HTTP 请求到服务器，并满足一定条件（例如合法的 token），
// 服务器会将该请求升级为 WebSocket 连接，建立双向通信通道。
// 需要：发送者id，接收者id，消息类型，发送的内容，发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并且校验token等的合法性
	//检验token
	query := request.URL.Query()
	//当前使用者
	Id := query.Get("userId")
	//10进制64位的数字
	userId, _ := strconv.ParseInt(Id, 10, 64)
	isvalida := true
	//升级为websocket连接
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Print(err)
		return
	}
	//2.获取conn
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登录时间
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
	}
	//3用户关系
	//4.将useid与node绑定并且加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送逻辑
	go sendProc(node)
	//6.完成接受逻辑
	go recvProc(node)

	//7.加如在线用户到缓存
	SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
}

// 发送消息逻辑，其实就是服务端接受一个消息
func sendProc(node *Node) {
	for {
		// 用于监听多个通道的消息，并执行相应的代码块。
		select {
		case data := <-node.DataQueue:
			fmt.Print("[ws]sendProc >>>> msg", string(data))
			// 通过 node.Conn 对象的 WriteMessage 方法向 WebSocket 连接发送消息。
			//websocket.TextMessage
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Print(err)
				return
			}
		}
	}
}

// 接受消息逻辑，其实是服务端发送一个消息，每个客户端都是独立享有一个该协程
func recvProc(node *Node) {
	for {
		//每次读取一条信息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Print(err)
			return
		}
		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Print(err)
			return
		}
		//如果是心跳检测的消息就将该用户的心跳时间更新
		//心跳检测 msg.Media==-1 || msg.Type==3
		if msg.Type == 3 { //进行心跳检测
			currenTime := uint64(time.Now().Unix())
			node.Heartbeat(currenTime)
		} else {
			//进行后端处理逻辑
			dispatch(data)
			fmt.Println("[ws] recvProc <<<< ", string(data))
		}
	}
}

// data为字节流的message，得将其转化为结构体Message/
func dispatch(data []byte) {
	msg := Message{}
	msg.CreateTime = uint64(time.Now().Unix())
	// 将字节数据 data 解析为一个 msg 对象，只会解析data中有的部分
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Print(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		fmt.Println("dispatch data:", string(data))
		sendMsg(msg.TargetId, data)
	case 2: //群发
		sendGroupMsg(msg.TargetId, data) //发送的群ID，消息内容
	}
}

func sendGroupMsg(TargetId int64, msg []byte) {
	fmt.Println("开始群发消息")
	userIds := SearchUserByGroupId(uint(TargetId))
	for i := 0; i < len(userIds); i++ {
		if TargetId != int64(userIds[i]) {
			sendMsg(int64(userIds[i]), msg)
		}
	}
}

// 将消息发送给userID，发送者存在msg.userid中
func sendMsg(userId int64, msg []byte) {
	//将userId给锁定到当前进程下
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()

	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	//用于管理 Redis 操作的超时、取消和其他相关信息。
	ctx := context.Background()
	//接受者
	targetIdStr := strconv.Itoa(int(userId))
	//发送者
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))

	//检查发送者的在线状态，是否掉线
	r, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Print("用户已经掉线了")
		fmt.Println(err)
		return
	}
	if r != "" {
		if ok { //发送了一条消息给这个人
			fmt.Println("sendMsg >>> userID: ", userId, "  msg:", string(msg))
			node.DataQueue <- msg
		}
	}

	//主要是这两个人之间总的聊天信息条数
	//这一块是用户获取聊天界面所有信息，即每次发送消息都会刷新用户聊天页面
	var key string
	//存入到redis中的键值
	if userId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	//用于获取 Redis 有序集合中指定范围内的成员列表
	// ，按照成员的分数从高到低排序。0,-1表示有序集的所有成员
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Print(err)
	}

	//将这条信息加入到redis中
	score := float64(cap(res)) + 1
	// &redis.Z{score, msg}：这是一个 redis.Z 类型的对象，
	// 表示要添加到有序集合的成员和分数。score 是分数，msg 是成员的值。
	ress, e := utils.Red.ZAdd(ctx, key, &redis.Z{Score: score, Member: msg}).Result()
	if e != nil {
		fmt.Print(e)
	}
	//通常情况下，它可能表示添加到有序集合中的成员数量、被更新的成员数量等
	fmt.Print(ress)
}

// 获取redis里面的消息，发送者-》接受者，redis中id在start到end中的所有数据，从小到大为0，从大到小为1
// 这一块是为了将用户信息刷新
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	//用于管理 Redis 操作的超时、取消和其他相关信息。
	ctx := context.Background()
	//发送者
	userIdStr := strconv.Itoa(int(userIdA))
	//接受者
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	var rels []string
	var err error
	//按照顺序取出这两个用户之间的所有聊天记录
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err)
	}

	//将这两个用户之间所有的信息返回
	return rels
}

//用户心跳机制

// 更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
}

// 清理掉超时的连接 ：传入一个对象返回他是否超时
func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanconnetion err", r)
		}
	}()
	// fmt.Print("定时任务，清理超时连接", param)
	// 枚举每个用户绑定的连接
	currenTime := uint64(time.Now().Unix())
	for i,_ := range clientMap {
		// fmt.Println(idx,":",currenTime,">>>>")
		node := clientMap[i]
		if node.IsHeartbeatTimeOut(currenTime) {
			fmt.Println("心跳超时。。。关闭连接：", node)
			//强制关闭连接
			node.Conn.Close()
		}
	}
	return result
}

// 检测用户心跳是否超时,对一个连接进行这个操作,传入一个当前时间
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Print("心跳超时。。。。自动下线", node)
		timeout = true
	}
	return
}
