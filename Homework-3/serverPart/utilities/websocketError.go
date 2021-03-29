package utilities

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

//排除因前端断开连接导致的错误，同时报告其他错误
func Error(coon *websocket.Conn, err error) {
	if err != nil {
		if websocket.IsCloseError(err, 1005) || strings.Contains(err.Error(),
			"An existing connection was forcibly closed by the remote host.") {
			fmt.Println("websocket: close")
		} else {
			coon.Close()
			log.Println(err)
		}
	}
}
