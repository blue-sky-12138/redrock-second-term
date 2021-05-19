package control

import (
	lb "SecondTerm/Homework-7/workloadManager/loadBalance"
	"io"
	"log"
	"net/http"
	"strings"
)

func RoutersEntrance() {
	e := &engine{}
	e.Run()
}

type engine struct{}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var url string
	if strings.HasPrefix(req.URL.Path, "/serve/oauth") || strings.HasPrefix(req.URL.Path, "/serve/user") {
		url = lb.UserBalance()
	} else if strings.HasPrefix(req.URL.Path, "/serve/video") {
		url = lb.VideoBalance()
	} else if strings.HasPrefix(req.URL.Path, "/serve/download") || strings.HasPrefix(req.URL.Path, "/serve/upload") {
		url = lb.FileBalance()
	}

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
	} else {
		quickForwarding(url, req, &w)
	}
}

func quickForwarding(url string, req *http.Request, w *http.ResponseWriter) {
	url += req.URL.String()

	sendReq, err := http.NewRequest(req.Method, url, req.Body)
	for key, val := range req.Header {
		sendReq.Header.Set(key, val[0])
	}

	resp, err := http.DefaultClient.Do(sendReq)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		(*w).Write([]byte("InternalServerError"))
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		(*w).Header().Set(k, v[0])
	}
	io.Copy((*w), resp.Body)
}

func (e *engine) Run(addr ...string) (err error) {
	defer func(err error) {
		if err != nil {
			log.Println(err)
		}
	}(err)

	var port string
	switch len(addr) {
	case 0:
		port = ":8080"
		log.Println("Listening address is not set. It is set to \":8080\"")
	case 1:
		port = addr[0]
		log.Println("The serve is set to \"" + port + "\"")
	default:
		panic("Listening address is too many")
	}

	err = http.ListenAndServe(port, e)
	return err
}
