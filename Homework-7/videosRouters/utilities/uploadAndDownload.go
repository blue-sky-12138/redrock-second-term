package utilities

var (
	HeadPath  = "./static/users/"
	VideoPath = "./static/videos/"
)

//获取视频文件及封面
type GetVideoFile struct {
	BvCode   string `uri:"bvCode"`   //视频bv号
	FileName string `uri:"fileName"` //文件名
}

//获取用户头像
type GetUserHead struct {
	Id       string `uri:"id"`       //用户id
	FileName string `uri:"fileName"` //文件名
}

//Attention!未完工状态
//获取投稿视频信息
type NewVideo struct {
	BvCode   string `form:"bv_code"`   //视频bv号
	Title    string `form:"title"`     //视频标题
	Brief    string `form:"brief"`     //简介
	DateTime string `form:"date_time"` //投稿时间
	AuthorId int64  `form:"user_id"`   //投稿人
	Joint    int    `form:"joint"`     //是否是联合投稿
	Type     int    `form:"type"`      //投稿类型(自制还是转载)
	Activity int    `form:"activity"`  //参加的活动

}
