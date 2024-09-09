// 初始化 与 app.yml关联
package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	// "ginchat/models" //用来测试数据
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	//配置文件名称
	viper.SetConfigName("app")
	//查找配置文件所在路径
	viper.AddConfigPath("config")
	//查找并读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited")
	// fmt.Println("config app", viper.Get("app"))
	// fmt.Println("config mysql", viper.Get("mysql"))
}
func InitMySQL() {
	//自定义日志打印打印SQL语句,前端页面点击执行sql语句时会有日志打印，便于后期调试
	newLogger := logger.New(
		//将日志输出到控制台
		// 这是 log.New 函数的参数，
		// 用于指定日志记录器的标志选项，log.LstdFlags 表示使用标准的日期和时间格式。
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值 表示执行时间超过该阈值的 SQL 查询将被记录为慢查询。
			LogLevel:      logger.Info, //日志级别，指定记录的日志消息级别
			Colorful:      true,        //是否启用彩色日志输出，设置为 true 表示启用
		},
	)
	//不要用冒号
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("mysql not inited")
	} else {
		fmt.Println("mysql inited")
	}
	//测试拿去数据
	// user := models.UserBasic{}
	// DB.Find(&user)
	// fmt.Println(user)
}

// addr: "127.0.0.1:6379"
// password: "20030204"
// DB: 0
// poolSize: 30
// minIdleconn: 30

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	ctx := context.Background()
	pong, err := Red.Ping(ctx).Result()
	if err != nil {
		fmt.Println("config redis inited 。。。。", err)
	} else {
		fmt.Println("Redis inited 。。。。", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
// ctx：上下文对象，可能包含与请求相关的信息。
// channel：要发布消息的Redis频道的名称。
// msg：要发布的消息内容。
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("publis.....", msg)
	err = Red.Publish(ctx, channel, msg).Err()

	return err
}

// Subscribe 订阅Redis指定频道中的消息
// ctx：上下文对象，可能包含与请求相关的信息。
// channel：要订阅的Redis频道的名称。
func Subscribe(ctx context.Context, channel string) (string, error) {
	//创建一个Redis订阅对象
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe....", ctx)
	//从订阅中接收消息。该函数将阻塞，直到接收到消息为止。
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("Subscribe....", msg.Payload)
	return msg.Payload, err
}
