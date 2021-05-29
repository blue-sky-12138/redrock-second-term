package BlueMQ

import (
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

//返回值为，连接是否生效
func Error(coon *websocket.Conn, err error) bool {
	if websocket.IsCloseError(err, 1005) || strings.Contains(err.Error(),
		"An existing connection was forcibly closed by the remote host.") {
		return true
	} else {
		coon.Close()
		log.Println(err)
		return false
	}
}
