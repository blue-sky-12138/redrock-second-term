package loadBalance

import (
	"sync"
)

type ServerInfo struct {
	IsWork bool
	Url    string
}

type ServerGroup struct {
	lock    sync.Mutex
	Turn    int
	Members []ServerInfo
}

var (
	UserServers  ServerGroup
	VideoServers ServerGroup
	FileServers  ServerGroup
)

func BalanceInit() {
	UserServers.Members = append([]ServerInfo{}, ServerInfo{
		IsWork: true,
		Url:    "http://localhost:8001",
	}, ServerInfo{
		IsWork: true,
		Url:    "http://localhost:8001",
	})

	VideoServers.Members = append([]ServerInfo{}, ServerInfo{
		IsWork: true,
		Url:    "http://localhost:8002",
	}, ServerInfo{
		IsWork: true,
		Url:    "http://localhost:8002",
	})

	FileServers.Members = append([]ServerInfo{}, ServerInfo{
		IsWork: true,
		Url:    "http://localhost:8003",
	}, ServerInfo{
		IsWork: true,
		Url:    "http://localhost:8003",
	})
}

func (s *ServerGroup) getServer() string {
	s.lock.Lock()
	defer s.lock.Unlock()

	originTurn := s.Turn
	for {
		s.Turn = (1 + len(s.Members)) % len(s.Members)

		if s.Members[s.Turn].IsWork {
			return s.Members[s.Turn].Url
		} else if s.Turn == originTurn {
			return ""
		}
	}
}

func UserBalance() string {
	url := UserServers.getServer()
	return url
}

func VideoBalance() string {
	url := VideoServers.getServer()
	return url
}

func FileBalance() string {
	url := FileServers.getServer()
	return url
}
