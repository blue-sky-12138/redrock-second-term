package main

import (
	hr "SecondTerm/Homework-5/HTTP_Request"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
)

const RequestUrl = "http://jwzx.cqupt.edu.cn/data/json_StudentSearch.php"

//理论上可以，但时间太长了……
//能够知道学号的组成，例如各学院的代号，可以排除一大部分无效学号，加快速度
func main() {
	number := 0
	for i := 2016000000; i < 2021000000; i++ {
		data := gin.H{
			"searchKey": i,
		}

		r, err := hr.Get(RequestUrl, data)
		if err != nil {
			log.Println(err)
		}

		regXM, _ := regexp.Compile(`"xm":"`)
		if !regXM.MatchString(r) {
			fmt.Println(i, "is false, continuing")
			continue
		} else {
			number++
		}
		//这里本来是想获取名字的，但获取的中文的转义字符不能转换成中文，暂时放弃
		//indexXM := regXM.FindAllStringIndex(r, 1)
		//
		//regEN, _ := regexp.Compile(`","xmEn"`)
		//indexEN := regEN.FindAllStringIndex(r, 1)
		//
		//Name := r[indexXM[0][1]:indexEN[0][0]]
		//Name = strings.ReplaceAll(Name, " ", "")
		//Name = fmt.Sprintf("%v", Name)
		//fmt.Println(Name)
	}
	fmt.Println(number)
}
