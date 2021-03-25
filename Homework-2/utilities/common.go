package utilities

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
)

var ErrorUserNotExist = errors.New("the user is not exist")

//用于与前端通信的响应统一结构
type Resp struct {
	Code    int         `json:"code"`    //响应代号
	Message string      `json:"message"` //响应信息
	Data    interface{} `json:"data"`    //返回数据
}

//通用信息
type Common struct {
	Date           string `json:"date"`           //日期
	Likes          int64  `json:"likes"`          //点赞数
	Coins          int64  `json:"coins"`          //投币数
	Collections    int64  `json:"collections"`    //收藏数
	Shares         int64  `json:"shares"`         //分享数
	CommentNumbers int64  `json:"comment_number"` //评论总数
}

//快捷返回数据
func (r *Resp) Return(ctx *gin.Context, code int, message string, data interface{}) {
	r.Code = code
	r.Message = message
	if data != nil {
		r.Data = screenData(&data)
	}

	ctx.JSON(200, r)
	return
}

//类似于Return()，不过这个方法不会筛选出零值，返回除Model(id保留)的所有数据
//如果Model的id值为零，也不会保留
func (r *Resp) ReturnAll(ctx *gin.Context, code int, message string, data interface{}) {
	r.Code = code
	r.Message = message
	if data != nil {
		r.Data = screenModel(&data)
	}

	ctx.JSON(200, r)
	return
}

//返回多个数据，要求输入时是key-value输入
func (r *Resp) ReturnMuti(ctx *gin.Context, code int, message string, args ...interface{}) {
	r.Code = code
	r.Message = message
	if l := len(args); args != nil {
		data := make(map[string]interface{})
		for i := 0; i < l; i += 2 {
			v := reflect.ValueOf(args[i])
			data[v.String()] = screenModel(&args[i+1])
		}
		r.Data = data
	}

	ctx.JSON(200, r)
	return
}

//筛选出Model里的除id外的变量
//Model要求一定是第一个字段
func screenModel(data *interface{}) interface{} {
	typ := reflect.TypeOf(*data)
	val := reflect.ValueOf(*data)

	if typ.Kind() == reflect.Slice { //如果是切片数据
		if val.IsZero() { //如果切片为空
			return nil
		} else {
			l := val.Len()
			child := make([]interface{}, 0)
			for i := 0; i < l; i++ { //分离切片递归
				v := val.Index(i).Interface()
				child = append(child, screenModel(&v))
			}

			return child
		}
	}

	result := make(map[string]interface{})
	fieldNum := typ.NumField()

	for i := 0; i < fieldNum; i++ {
		t := typ.Field(i)
		v := val.Field(i)

		if t.Name == "Model" { //处理Model数据
			if !v.Field(0).IsZero() { //如果id值不是零值
				result["id"] = v.Field(0).Interface()
			}
			continue
		}

		if v.Kind() == reflect.Struct || v.Kind() == reflect.Slice { //如果变量类型是结构体或切片，递归
			str := v.Interface()
			result[t.Tag.Get("json")] = screenData(&str)
			continue
		}

		if t.Tag.Get("json") == "" {
			continue
		}

		//以json标签作为key来加入到map中
		result[t.Tag.Get("json")] = v.Interface()
	}

	return result
}

//筛选出所有非零值的变量
func screenData(data *interface{}) interface{} {
	typ := reflect.TypeOf(*data)
	val := reflect.ValueOf(*data)

	if typ.Kind() == reflect.Slice { //如果是切片数据
		if val.IsZero() { //如果切片为空
			return nil
		} else {
			l := val.Len()
			child := make([]interface{}, 0)
			for i := 0; i < l; i++ { //分离切片递归
				v := val.Index(i).Interface()
				child = append(child, screenData(&v))
			}

			return child
		}
	}

	result := make(map[string]interface{})
	fieldNum := typ.NumField()

	for i := 0; i < fieldNum; i++ {
		t := typ.Field(i)
		v := val.Field(i)

		if t.Name == "Model" || v.IsZero() { //忽略零值的变量和gorm.model
			continue
		}

		if v.Kind() == reflect.Struct || v.Kind() == reflect.Slice { //如果变量类型是结构体或切片，递归
			str := v.Interface()
			result[t.Tag.Get("json")] = screenData(&str)
			continue
		}

		if t.Tag.Get("json") == "" {
			continue
		}

		//以json标签作为key来加入到map中
		result[t.Tag.Get("json")] = v.Interface()
	}

	return result
}
