package api

import (
	"AirGo/global"
	"AirGo/service"
	"AirGo/utils/response"
	"AirGo/utils/websocket_plugin"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

// websocket im 测试
func WebSocketMsg(ctx *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {

			return true
		},
		//后端带token响应，否则前端接收不到数据
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	uID, ok := ctx.Get("uID")
	if !ok || uID == nil {
		global.Logrus.Error("uID参数错误")
		response.Fail("uID参数错误", nil, ctx)
		return
	}
	uIDInt := uID.(int)

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//ebsocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header"
		//nginx:
		//proxy_set_header Upgrade $http_upgrade;
		//proxy_set_header Connection upgrade;
		//proxy_set_header X-Real-IP $remote_addr;
		global.Logrus.Error("websocket upgrade error:", err)
		response.Fail("websocket err:"+err.Error(), nil, ctx)
		return
	}
	//defer conn.Close()
	client := &websocket_plugin.Client{
		ID: strconv.Itoa(uIDInt),
		//ID:            ctx.ClientIP(),
		WsSocket:      conn,
		ClientChannel: make(chan []byte),
		ExpireTime:    5 * time.Second, //5秒过期时间
		QuitChanel:    make(chan bool),
	}
	global.WsManager.OnlineChannel <- client
	go client.Read(global.WsManager, service.GetNodesStatus)
	go client.Write(global.WsManager)
}
