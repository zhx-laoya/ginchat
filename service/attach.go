// // 上传相关文件
package service

import (
	"fmt"
	"ginchat/utils"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	UploadLocal(c)
}

// 上传文件到本地 tag未注释
func UploadLocal(c *gin.Context) {
	// 将 c.Writer 赋值给变量 w，用于向客户端发送响应
	w := c.Writer
	// 用于获取客户端的请求信息。
	req := c.Request
	// 使用 req.FormFile 方法获取客户端上传的文件。
	// 它通过表单字段名 "file" 来获取文件对象，
	// 并返回文件对象 srcFile、文件头信息 head，以及可能发生的错误 err。
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	// 设置一个默认的文件后缀为 .png。
	suffix := ".png"
	// 从文件头信息 head 中获取原始文件名，并赋值给变量 ofilName。
	ofilName := head.Filename
	// 使用 strings.Split 方法按照文件名中的 . 进行分割，
	// 将文件名拆分为多个部分，并将结果赋值给变量 tem。
	tem := strings.Split(ofilName, ".")
	// 如果文件名中包含 .
	// 则将最后一个部分作为文件后缀，并将其赋值给变量 suffix。
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	// 根据当前时间戳、随机数和文件后缀生成一个唯一的文件名
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)

	// os.Create 方法创建一个新的文件，路径为 ./asset/upload/ 加上生成的文件名 fileName，
	// 并将文件对象赋值给变量 dstFile。
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	//copy是为了后续能够找到并传给用户观看
	// io.Copy 方法将客户端上传的文件内容复制到创建的目标文件中。
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	url := "./asset/upload/" + fileName
	utils.RespOK(w, url, "发送图片成功")
}
