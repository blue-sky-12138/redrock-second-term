package main

import (
	hr "SecondTerm/Homework-5/HTTP_Request"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type classInfo struct {
	name      string
	position  string
	week      [][]int
	time      int
	teacher   string
	classType string
}

type schedule struct {
	nowWeek   int
	classList [6][7][]classInfo
}

const (
	classScheduleUrl = "http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php"
)

const (
	time1and2 = iota
	time3and4
	time5and6
	time7and8
	time9and10
	timeThreeClass
)

func main() {
	menu()
}

func (res *schedule) GetSchedule(id int) error {
	var data = gin.H{
		"xh": id,
	}

	r, err := hr.Get(classScheduleUrl, data)
	if err != nil {
		return err
	}

	err = res.AnalysisSchedule(r)
	if err != nil {
		return err
	}

	return nil
}

func (res *schedule) AnalysisSchedule(data string) error {
	var (
		temClassData [6][7]string
		temTr        []string
		err          error
	)
	//提取出表格形式的数据
	reg := regexp.MustCompile(`<div id="kbStuTabs-table">[\s\S]*<div id="kbStuTabs-list">`)
	match := reg.FindAllString(data, -1)
	if len(match) != 1 {
		return fmt.Errorf("data error")
	}
	s := match[0]

	//获取第几周
	reg = regexp.MustCompile(`当前周次：</b>第(\S*)周`)
	res.nowWeek, err = strconv.Atoi(reg.ReplaceAllString(reg.FindAllString(s, -1)[0], "$1"))
	if err != nil {
		return fmt.Errorf("nowWeek error")
	}

	//分离出每一行
	reg = regexp.MustCompile(`<tr`)
	index := reg.FindAllStringIndex(s, -1)
	for i := range index {
		if i == 1 || i == 2 || i == 4 || i == 5 || i == 7 || i == 8 {
			temTr = append(temTr, s[index[i][1]:index[i+1][0]])
		}
	}

	//分离出每一个单元格
	for i, v := range temTr {
		reg = regexp.MustCompile(`<td`)
		index = reg.FindAllStringIndex(v, -1)
		for j := range index {
			if j+1 == len(index) {
				temClassData[i][j-1] = v[index[j][1]:]
			} else if j > 0 {
				temClassData[i][j-1] = v[index[j][1]:index[j+1][0]]
			}
		}
	}

	//加工每一节课的信息
	var (
		temClassinfo classInfo
		infoSlice    []classInfo
	)

	for i := range temClassData {
		for j, v := range temClassData[i] {
			//获得上课的具体时间
			var time int
			reg = regexp.MustCompile(`3节连上`)
			if reg.MatchString(v) {
				time = timeThreeClass
			} else {
				time = i
			}

			//获取课程名，地点，上课周数
			reg1 := regexp.MustCompile(`[A-Z]\d{7}-(?P<name>\S*)<br>地点：(?P<position>\S*)\s*<br>(?P<week>\S*)<`)
			match1 := reg1.FindAllString(v, -1)

			//获得老师和选课类型
			reg2 := regexp.MustCompile(`>(?P<teacher>\S*)\s*(?P<classType>\S修)\s`)
			match2 := reg2.FindAllString(v, -1)

			if len(match1) != len(match2) {
				return fmt.Errorf("classInfo error")
			} else if match1 == nil {
				continue
			} else {
				//匹配数值
				l := len(match1)
				for k := 0; k < l; k++ {
					temClassinfo.name = reg1.ReplaceAllString(match1[k], "$name")
					temClassinfo.position = reg1.ReplaceAllString(match1[k], "$position")
					temClassinfo.analysisWeek(reg1.ReplaceAllString(match1[k], "$week"))

					temClassinfo.teacher = reg2.ReplaceAllString(match2[k], "$teacher")
					temClassinfo.classType = reg2.ReplaceAllString(match2[k], "$classType")

					temClassinfo.time = time

					//将结果加入临时切片中
					infoSlice = append(infoSlice, temClassinfo)
					//重置临时变量的数值
					temClassinfo = classInfo{}
				}
			}

			//赋予结果
			res.classList[i][j] = append(res.classList[i][j], infoSlice...)

			//重置临时变量的数值
			infoSlice = nil
		}
	}

	return nil
}

func (sed *classInfo) analysisWeek(data string) {
	var tem []int

	//除去中文
	data = strings.ReplaceAll(data, "周", "")

	//将多个时间段分隔
	s1 := strings.Split(data, ",")
	for _, v := range s1 {
		//将时间段分隔
		s2 := strings.Split(v, "-")
		if len(s2) == 2 {
			m, _ := strconv.Atoi(s2[0])
			n, _ := strconv.Atoi(s2[1])
			tem = append(tem, m, n)

			sed.week = append(sed.week, tem)
		} else {
			m, _ := strconv.Atoi(s2[0])
			tem = append(tem, m, 0)

			sed.week = append(sed.week, tem)
		}

		//重置临时时间值
		tem = []int{}
	}
}

//这是formatPrint的快捷方式，自动打印当前周的课表
func (sed *schedule) shouldFormatPrint() {
	sed.formatPrint(sed.nowWeek)
}

func (sed *schedule) formatPrint(week int) {
	//先遍历列，再遍历行
	for j := 0; j < 7; j++ {
		switch j {
		case 0:
			fmt.Println("星期一")
		case 1:
			fmt.Println("星期二")
		case 2:
			fmt.Println("星期三")
		case 3:
			fmt.Println("星期四")
		case 4:
			fmt.Println("星期五")
		case 5:
			fmt.Println("星期六")
		case 6:
			fmt.Println("星期日")
		}

		for i := 0; i < 6; i++ {
			//取出该天的全部课程
			for _, v := range sed.classList[i][j] {
				v.formatPrint(week)
			}
		}

		fmt.Printf("\n")
	}
}

func (ci *classInfo) formatPrint(week int) {
	//筛选
	if len(ci.week) == 0 {
		return
	} else {
		ok := 0
		for _, v := range ci.week {
			if v[1] == 0 && week == v[0] {
				ok += 1
			} else if v[1] != 0 && week >= v[0] && week <= v[1] {
				ok += 1
			}
		}

		//如果没有一个时间是符合的，直接退出
		if ok == 0 {
			return
		}
	}

	switch ci.time {
	case time1and2:
		fmt.Printf("7:00-8:40   ")
	case time3and4:
		fmt.Printf("10:15-11:55 ")
	case time5and6:
		fmt.Printf("14:00-15:40 ")
	case time7and8:
		fmt.Printf("16:15-17:55 ")
	case time9and10:
		fmt.Printf("19:00-20:40 ")
	case timeThreeClass:
		fmt.Printf("19:00-21:35 ")
	}

	fmt.Println(ci.name, ci.position, ci.teacher, ci.classType)
}

func menu() {
	var (
		id   int
		week int
	)

	for {
		fmt.Printf("请输入查询的id:")
		n, err := fmt.Scanln(&id)
		if err != nil || n != 1 {
			fmt.Println("data error")
			continue
		}

		fmt.Printf("请输入查询的周数(0为当前周):")
		n, err = fmt.Scanln(&week)
		if err != nil || n != 1 {
			fmt.Println("data error")
			continue
		} else {
			break
		}
	}

	var sc schedule
	err := sc.GetSchedule(id)
	if err != nil {
		log.Println(err)
		return
	}

	if week == 0 {
		sc.shouldFormatPrint()
	} else {
		sc.formatPrint(week)
	}
}
