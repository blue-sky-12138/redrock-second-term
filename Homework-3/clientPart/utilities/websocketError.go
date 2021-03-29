package utilities

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

func Error(coon *websocket.Conn, err error) {
	if err != nil {
		if websocket.IsCloseError(err, 1005) || strings.Contains(err.Error(),
			"An existing connection was forcibly closed by the remote host.") {
			fmt.Println("与服务器的连接已断开，稍后该进程将结束")
		} else {
			coon.Close()
			log.Println(err)
		}
	}
}
