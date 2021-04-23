package HTTP_Request

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

func Get(url string, data interface{}) (result string, err error) {
	d := fillUrlValue(&data)
	if d != "" {
		url += "?" + d
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

//自动填充form参数
func fillUrlValue(data *interface{}) string {
	//排除空数据
	if (*data) == nil {
		return ""
	}

	typ := reflect.TypeOf(*data)
	val := reflect.ValueOf(*data)
	res := ""

	//兼容map
	if val.Kind() == reflect.Map {
		for _, v := range val.MapKeys() {
			d := reflect.ValueOf(val.MapIndex(v).Interface())
			switch d.Kind() {
			case reflect.String:
				res += v.String() + "=" + d.String() + "&"
			case reflect.Int:
				res += v.String() + "=" + strconv.FormatInt(d.Int(), 10) + "&"
			}

		}

		if len(res) < 1 {
			return ""
		} else {
			return res[:len(res)-1]
		}
	}

	fieldNum := typ.NumField()

	for i := 0; i < fieldNum; i++ {
		t := typ.Field(i)
		v := val.Field(i)

		//以form标签作为key来加入到values中
		if v.Kind() == reflect.Int {
			res += t.Tag.Get("form") + "=" + strconv.FormatInt(v.Int(), 10) + "&"
		} else {
			res += t.Tag.Get("form") + "=" + v.String() + "&"
		}
	}

	if len(res) < 1 {
		return ""
	} else {
		return res[:len(res)-1]
	}
}
