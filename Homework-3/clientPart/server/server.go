package server

type Room struct {
	ID      int `json:"id"`
	Numbers int `json:"numbers"`
}

type Rooms struct {
	RS []Room `json:"data"`
}

type Message struct {
	Type    int    `json:"type"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	Content string `json:"content"`
}
