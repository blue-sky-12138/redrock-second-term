package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"serverPart/model"
	ut "serverPart/utilities"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	largestNumber = 200 //聊天室最大人数
)

type Manage struct {
	Rooms map[int]Room
}

type Room struct {
	ID      int `json:"id"`
	Numbers int `json:"numbers"`
	Clients map[string]*Client
	Filters ut.Filters
}

type Message struct {
	Type    int    `json:"type"`    //类型，每个客户端首个信息传输操作类型
	Name    string `json:"name"`    //用户名
	Time    string `json:"time"`    //时间
	Content string `json:"content"` //信息内容，每个客户端首个信息传输房间号
}

var (
	manageLock sync.Mutex
	manage     = Manage{Rooms: make(map[int]Room)}
	filters    ut.Filters
	banners    = make(map[string][]string)
)

func WebsocketConnect(ctx *gin.Context) {
	up := websocket.Upgrader{
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	coon, err := up.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("connectError:", err)
		return
	}

	c := NewClient(coon)
	go c.Receive()
	go c.Send()

	//播报用户进入房间
	msg := Message{
		Type:    100,
		Name:    "Room(" + strconv.Itoa(manage.Rooms[c.RoomID].Numbers) + "人)",
		Content: c.Name + "进入了房间",
	}
	manage.Rooms[c.RoomID].Broadcast(msg)

	return
}

func List(ctx *gin.Context) {
	var (
		resp   ut.Resp
		result []Room
	)
	for _, v := range manage.Rooms {
		result = append(result, v)
	}
	resp.ReturnAll(ctx, 400, "", result)
}

func (r Room) Broadcast(msg Message) {
	if msg.Type == 200 { //如果是纯信息，进行过滤
		msg.Filter()
	}

	for _, c := range r.Clients {
		if msg.BanOrNot(c) {
			c.MsgChan <- msg
		}
	}
}

func (m Manage) Broadcast(msg Message) {
	manageLock.Lock()
	for _, v := range manage.Rooms {
		v.Broadcast(msg)
	}
	manageLock.Unlock()
}

func (m *Manage) EnterRoom(c *Client) {
	manageLock.Lock()

	r := m.Rooms[c.RoomID]
	r.Numbers++
	r.Clients[c.Name] = c

	m.Rooms[c.RoomID] = r

	manageLock.Unlock()
}

func (m *Manage) QuitRoom(c *Client) {
	manageLock.Lock()

	//更新房间信息
	r := m.Rooms[c.RoomID]
	r.Numbers--

	//删除房间中客户端资源
	delete(r.Clients, c.Name)

	//如果房间里没有人，删除房间
	if r.Numbers <= 0 {
		delete(m.Rooms, c.RoomID)
	} else { //否则，更新管理员的房间信息
		m.Rooms[c.RoomID] = r

		//播报用户离开房间
		msg := Message{
			Name:    "Room(" + strconv.Itoa(manage.Rooms[c.RoomID].Numbers) + "人)",
			Content: c.Name + "离开了房间",
		}
		m.Rooms[c.RoomID].Broadcast(msg)
	}

	manageLock.Unlock()
}

func (m *Manage) CreateRoom(c *Client) {
	manageLock.Lock()

	//判断是否存在
	_, ok := m.Rooms[c.RoomID]
	if ok { //如果存在该房间，加入
		manageLock.Unlock()
		m.EnterRoom(c)
		return
	}

	r := Room{
		ID:      c.RoomID,
		Numbers: 1,
		Clients: make(map[string]*Client),
	}
	r.Clients[c.Name] = c

	m.Rooms[c.RoomID] = r
	manageLock.Unlock()
}

//敏感词过滤
func (msg *Message) Filter() {
	for _, v := range filters {
		msg.Content = strings.Replace(msg.Content, v.Content, "***", -1)
	}
}

//更新过滤器
func FilterUpdate() {
	for {
		model.FiltersUpdate(&filters)
		time.Sleep(1 * time.Hour)
	}
}

func (msg *Message) BanOrNot(c *Client) bool {
	bans := banners[c.Name]
	for _, v := range bans {
		if v == msg.Name { //如果接收方的黑名单里有发信息的人
			return false
		}
	}
	return true
}

//更新屏蔽器
func BannerUpdate() {
	for {
		model.BannersUpdate(&banners)
		time.Sleep(1 * time.Hour)
	}
}
