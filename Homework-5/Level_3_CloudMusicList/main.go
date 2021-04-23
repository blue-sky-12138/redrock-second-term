package main

import (
	hr "SecondTerm/Homework-5/HTTP_Request"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
)

const CouldMusicListUrl = "https://music.163.com/discover/playlist/"

type ListInfo struct {
	Name   string
	Author string
	Play   string
	Url    string
}

//本来是想做获取网易云歌单的，但是调试中途golang对目标URL进行访问时，返回总是为"服务器开小差"，但在postman上又返回正常
//目前无法测试，只有大概的框架
func main() {
	i := 0
	var li []ListInfo
	for {
		var data = gin.H{
			"order":  "hot",
			"limit":  35,
			"offset": i * 35,
		}
		r, err := hr.Get(CouldMusicListUrl, data)
		if err != nil {
			log.Println(err)
			return
		}

		lice, err := AnalysisData(r)
		if err != nil {
			log.Println(err)
			return
		} else {
			li = append(li, lice...)
		}

		//只搜索6页
		if i == 5 {
			break
		} else {
			i++
		}
	}

	for _, v := range li {
		v.formatPrint()
	}
}

func (li *ListInfo) formatPrint() {
	fmt.Println("歌单名:", li.Name, "\t作者:", li.Author, "\t播放量:", li.Play, "\tURL:", li.Url)
}

func AnalysisData(data string) ([]ListInfo, error) {
	var (
		temlice []string
		res     []ListInfo
	)

	reg := regexp.MustCompile(`很抱歉，服务器开小差了，请稍后再试`)
	if reg.MatchString(data) {
		return nil, fmt.Errorf("the website is busy")
	}

	//分理出歌单列表
	reg = regexp.MustCompile(`<ul class="m-cvrlst f-cb" id="m-pl-container">(?P<content>[\s\S]*)<div id="m-pl-pager">`)
	index := reg.FindAllStringIndex(data, -1)
	if len(index) != 1 {
		return nil, fmt.Errorf("data error")
	}
	data = data[index[0][0]:index[0][1]]

	//分理出每一个歌单
	reg = regexp.MustCompile(`<li>`)
	index = reg.FindAllStringIndex(data, -1)
	for i := range index {
		if i+1 != len(index) {
			temlice = append(temlice, data[index[i][1]:index[i+1][0]])
		} else {
			temlice = append(temlice, data[index[i][1]:])
		}
	}

	//提取每个歌单的信息
	reg = regexp.MustCompile(`<span class="nb">(?P<play>[\S]*)</s[\S\s]*title="(?P<name>[\S\s]*)\shref="(?P<url>[\S]*)"\sclass[\s\S]*<a title="(?P<author>[\S\s]*)" href="/user/home`)
	for _, v := range temlice {
		sLice := reg.FindAllString(v, -1)
		if len(sLice) != 1 {
			return nil, fmt.Errorf("ListContent Error")
		}

		res = append(res, ListInfo{
			Name:   reg.ReplaceAllString(sLice[0], "$name"),
			Author: reg.ReplaceAllString(sLice[0], "$author"),
			Play:   reg.ReplaceAllString(sLice[0], "$play"),
			Url:    reg.ReplaceAllString(sLice[0], "$url"),
		})
	}

	return res, nil
}
