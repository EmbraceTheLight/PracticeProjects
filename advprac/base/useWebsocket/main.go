// websocket 基本使用
package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
			return
		}
		defer conn.Close()

		fmt.Println("Websocket 连接已建立")
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("读取消息失败:", err)
				return
			}
			fmt.Printf("收到消息: %s\n", message)
			message = append(message, []byte("已收到")...)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				fmt.Println("消息写入失败:", err)
				break
			}
		}

		fmt.Println("websocket 连接已关闭")
	})

	http.ListenAndServe(":8080", nil)
}
