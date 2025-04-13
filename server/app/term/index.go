package term

import (
	"errors"
	"net/http"
	"runtime"
	"server/core/app/request"
	"server/core/app/response"
	"server/middleware"
	"server/utils/cache"
	"server/utils/config"
	"server/utils/data"
	"server/utils/logger"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	pty2 "github.com/aymanbagabas/go-pty"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// TremApp 终端应用结构体
type TremApp struct {
}

// AddRoutes 添加终端相关路由
// parentGroup: 父路由组
func AddRoutes(parentGroup *gin.RouterGroup) {
	app := TremApp{}
	group1 := parentGroup.Group("/auth", middleware.AuthMiddleware) // 需要认证的路由组
	group2 := parentGroup.Group("/ws")                              // WebSocket路由组
	group1.GET("/token/:termId", app.GetToken)                      // 获取终端token
	group2.GET("/:clientId/:token", app.OpenWs)                     // 打开WebSocket连接
}

// clients 存储所有客户端连接的映射
var clients = make(map[string]*websocket.Conn)

// clientMutex 用于保护clients映射的互斥锁
var clientMutex sync.Mutex

// MsgType 消息类型枚举
type MsgType int

// 定义不同的消息类型常量
const (
	HEARTBEAT MsgType = iota // 心跳消息
	CONN                     // 连接消息
	CLOSE                    // 关闭消息
	MSG                      // 普通消息
	RESIZE                   // 调整大小消息
	ERROR                    // 错误消息
)

// TermData 终端数据结构体
type TermData struct {
	Type MsgType `json:"type"` // 消息类型
	Data string  `json:"data"` // 消息内容
}

// ToString 将终端数据转换为字符串
func (td *TermData) ToString() string {
	return convertor.ToString(td.Type) + ":" + td.Data
}

// Bytes 将终端数据转换为字节数组
func (td *TermData) Bytes() []byte {
	return []byte(td.ToString())
}

// Int 获取消息类型的整数值
func (td *TermData) Int() int {
	return int(td.Type)
}

// GetToken 获取终端连接的令牌
// 根据终端ID和用户ID生成一个临时令牌供WebSocket连接使用
func (app TremApp) GetToken(ctx *gin.Context) {
	termId := ctx.Param("termId")
	userId := convertor.ToString(request.GetUserID(ctx))
	token, _ := random.UUIdV4()
	clientId := "term_" + termId + "_" + userId
	// 设置令牌缓存，有效期30秒
	cache.GetCacheSystem().SetExpire(clientId, token, 30)
	response.Data(ctx, "", data.Map{"clientId": clientId, "token": token})
}

// CheckToken 验证终端连接令牌是否有效
// clientId: 客户端ID
// token: 客户端提供的令牌
// 返回: 令牌是否有效
func CheckToken(clientId, token string) bool {
	realToken, err := cache.GetCacheSystem().Get(clientId)
	if err != nil {
		logger.LOG.Error(err)
		return false
	}
	isValied := err == nil && token == realToken
	if isValied {
		// 令牌有效则删除缓存，确保一次性使用
		cache.GetCacheSystem().Delete(clientId)
	}
	return isValied
}

// OpenWs 打开WebSocket连接
// 验证令牌后建立WebSocket连接，并处理消息交互
func (app TremApp) OpenWs(ctx *gin.Context) {
	clientId := ctx.Param("clientId")
	token := ctx.Param("token")
	isValied := CheckToken(clientId, token)
	logger.LOG.Infof("clientId:%s, token:%s ,isValied:%v\n", clientId, token, isValied)
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 根据判断token的方法来鉴权,如果没token就返回false
			return isValied
		},
	}
	//升级get请求为webSocket协议
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.LOG.Errorln("websocket连接错误:" + err.Error())
		return
	}
	// 加入新客户端
	clientMutex.Lock()
	clients[clientId] = conn // 将新客户端添加到客户端集合中
	clientMutex.Unlock()
	// 处理WebSocket消息
	go handleWebSocketMessages(clientId, conn)
}

// sendMessage 发送消息到WebSocket客户端
// conn: WebSocket连接
// td: 要发送的终端数据
func sendMessage(conn *websocket.Conn, td TermData) {
	// 添加错误处理
	if conn == nil {
		return
	}

	// 使用写入锁保护WebSocket写入
	if err := conn.WriteMessage(websocket.TextMessage, td.Bytes()); err != nil {
		logger.LOG.Errorf("发送WebSocket消息错误: %s", err.Error())
	}
}

// closeWs 关闭WebSocket连接
// clientId: 客户端ID
// conn: WebSocket连接
func closeWs(clientId string, conn *websocket.Conn) {
	logger.LOG.Infof("closeWs clientId:%s\n", clientId)
	clientMutex.Lock()
	defer clientMutex.Unlock()
	delete(clients, clientId) // 删除客户端连接
	conn.Close()              // 关闭 WebSocket 连接
}

// sendError 发送错误消息到WebSocket客户端
// conn: WebSocket连接
// err: 要发送的错误信息
func sendError(conn *websocket.Conn, err error) {
	td := TermData{}
	td.Type = ERROR
	td.Data = err.Error()
	sendMessage(conn, td)
}

// handleWebSocketMessages 处理WebSocket消息
// clientId: 客户端ID
// conn: WebSocket连接
func handleWebSocketMessages(clientId string, conn *websocket.Conn) {
	flag := true
	// 启动pty终端
	ptmx, err := pty2.New()
	if err != nil {
		logger.LOG.Errorln("启动pty终端失败:" + err.Error())
		sendError(conn, err)
		closeWs(clientId, conn)
		return
	}
	defer ptmx.Close()
	var c *pty2.Cmd
	// 根据不同操作系统启动不同的终端
	if runtime.GOOS == "windows" {
		// Windows系统优先尝试启动PowerShell，失败则尝试启动cmd
		c = ptmx.Command(`powershell`)
		if err := c.Start(); err != nil {
			msg := err.Error()
			c = ptmx.Command(`cmd`)
			if err := c.Start(); err != nil {
				logger.LOG.Errorln("启动pty终端失败:" + msg)
				sendError(conn, errors.New(msg))
				closeWs(clientId, conn)
				return
			}
		}
	} else {
		// 非Windows系统使用配置的命令，失败则尝试启动sh
		c = ptmx.Command(config.CONF.Term.Command)
		if err := c.Start(); err != nil {
			msg := err.Error()
			c = ptmx.Command(`sh`)
			if err := c.Start(); err != nil {
				logger.LOG.Errorln("启动pty终端失败:" + msg)
				sendError(conn, errors.New(msg))
				closeWs(clientId, conn)
				return
			}
		}
	}
	// 读取终端响应发送给ws客户端
	go func() {
		buf := make([]byte, 512)
		var payload, overflow []byte
		for flag {
			// 从终端读取数据
			n, err := ptmx.Read(buf)
			if err != nil {
				logger.LOG.Errorln("读取pty终端数据错误:" + err.Error())
				sendError(conn, err)
				flag = false
				closeWs(clientId, conn)
				return
			}
			payload = append(payload[0:], overflow...)
			overflow = nil
			payload = append(payload, buf[:n]...)
			// 确保发送的是有效的UTF-8编码
			for !utf8.Valid(payload) {
				overflow = append(overflow[:0], append(payload[len(payload)-1:], overflow[0:]...)...)
				payload = payload[:len(payload)-1]
			}
			if len(payload) >= 1 {
				td := TermData{}
				td.Type = MSG
				td.Data = string(payload)
				sendMessage(conn, td)
			}
			payload = nil
		}
	}()
	// 读取ws客户端数据发送到pty终端
	for flag {
		_, p, err := conn.ReadMessage() // 阻塞等待读取客户端发送的消息
		if err != nil {
			logger.LOG.Errorln("读取WS数据错误:" + err.Error())
			flag = false
			break // 无法读取信息，退出循环
		}
		msg := string(p)
		mi := strings.Index(msg, ":")
		// 无冒号 无效数据
		if mi < 1 {
			continue
		}
		// 解析消息类型
		t, err := strconv.Atoi(msg[0:mi])
		// 类型错误 无效数据
		if err != nil {
			continue
		}
		// 提取消息数据
		data := msg[mi+1:]
		switch MsgType(t) {
		case HEARTBEAT:
			// 心跳消息处理，返回当前时间戳
			td := TermData{}
			td.Type = HEARTBEAT
			td.Data = convertor.ToString(time.Now().UnixMilli())
			sendMessage(conn, td)
		case CLOSE:
			// 收到客户端关闭消息，关闭服务端连接
			flag = false
			closeWs(clientId, conn)
			return
		case MSG:
			// 收到输入字符，发送到终端
			ptmx.Write([]byte(data))
		case RESIZE:
			// 收到终端调整大小请求
			size := strings.Split(data, ",")
			cols, _ := convertor.ToInt(size[0])
			rows, _ := convertor.ToInt(size[1])
			if cols > 0 && rows > 0 {
				// 调整终端窗口大小
				ptmx.Resize(int(cols), int(rows))
			}
		case ERROR:
			// 收到客户端错误消息，记录错误并关闭连接
			flag = false
			closeWs(clientId, conn)
			return
		default:
			// 未知消息类型，忽略
		}

	}
	closeWs(clientId, conn)
}
