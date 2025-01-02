package go_ws_sh

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"
)

func createhandleWebSocket(session Session) func(w context.Context, r *app.RequestContext) {
	return func(w context.Context, r *app.RequestContext) {
		codec, err := create_msg_codec()
		if err != nil {
			log.Println(err)
			// r.SetStatusCode(http.StatusInternalServerError)
			r.AbortWithMsg("Internal Server Error"+"\n"+err.Error(), consts.StatusInternalServerError)
			// w.WriteHeader(http.StatusInternalServerError)
			// w.Write([]byte("Internal Server Error" + "\n" + err.Error()))
			return
		}

		upgrader := websocket.HertzUpgrader{
			CheckOrigin: func(ctx *app.RequestContext) bool {
				return true // 允许跨域
			},
			EnableCompression: true,
		}
		err = upgrader.Upgrade(r, func(conn *websocket.Conn) {
			err := HandleWebSocketProcess(session, codec, conn)
			if err != nil {
				log.Println(err)
			}

		})
		// 升级HTTP连接到WebSocket协议
		// conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			r.AbortWithMsg("Internal Server Error"+"\n"+err.Error(), consts.StatusInternalServerError)
			// w.WriteHeader(http.StatusInternalServerError)
			// w.Write([]byte("Internal Server Error" + "\n" + err.Error()))
			return
		}
		// 设置标准输入、输出和错误流
		// 启动命令
		// 处理标准输出和错误流
		// 等待命令执行完毕
		//进程结束
		// 处理标准输入流
		//读取out_queue,并转换
		//将数据转换为二进制
		// 发送消息到WebSocket连接
		//读取out_queue,并转换
		//将数据转换为二进制
		// 发送消息到WebSocket连接
		// 循环读取WebSocket消息
		//连接结束
		// 将消息发送回客户端
		// err = conn.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }

	}
}
