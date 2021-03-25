package utilities

import "gorm.io/gorm"

//查询视频id的vpi
type VideoAPI struct {
	gorm.Model
	BVCode string `json:"bv_code"`
	P      int    `json:"p"`
}

//楼中楼评论结构
type ReplyComment struct {
	gorm.Model
	Author        UserInformation `json:"author" gorm:""` //评论作者信息
	AuthorId      int             `gorm:"type:int unsigned"`
	Content       string          `json:"content"` //评论内容
	DateTime      string          `json:"date_time"`
	Likes         int             `json:"likes"  gorm:"-"`
	ReplyAuthor   UserInformation `json:"reply_author" gorm:""` //被回复人的信息
	ReplyAuthorId int             `gorm:"type:int unsigned"`
}

//元评论结构构
type MetaComment struct {
	gorm.Model
	Author       UserInformation `json:"author" gorm:""`
	AuthorId     int             `gorm:"type:int unsigned"`
	Content      string          `json:"content"`
	Floor        int             `json:"floor"` //几楼
	DateTime     string          `json:"date_time"`
	Likes        int             `json:"likes" gorm:"-"`
	Heat         int             `json:"heat" gorm:"-"`                                              //评论热度，数值为点赞与踩的数值和
	ReplyComment []ReplyComment  `json:"comments_in_floor" gorm:"-;joinReferences:reply_comment_id"` //该楼的楼中楼评论
}

//视频新增评论结构
type NewComment struct {
	gorm.Model
	VideoId        int64  `json:"video_id"`         //视频ID
	DateTime       string `json:"date_time"`        //评论时间
	Content        string `json:"content"`          //评论内容
	AuthorId       int64  `json:"user_id"`          //评论者ID
	Floor          int    `json:"floor"`            //元评论楼层
	ReplyCommentId int64  `json:"reply_comment_id"` //被评论id
	ReplyAuthorId  int64  `json:"reply_author_id"`  //被评论者id
}

//视频操作信息结构
type OperateVideoInformation struct {
	gorm.Model
	UserId  int64 //用于存储已登录的用户ID
	VideoId int64
	Like    int `gorm:"column:likes"`       //点赞情况
	Coin    int `gorm:"column:coins"`       //投币情况
	Collect int `gorm:"column:collections"` //收藏情况
	Share   int `gorm:"column:shares"`      //分享情况
}

//视频作者信息结构
type VideoAuthorInformation struct {
	gorm.Model
	Name      string `json:"name"`      //用户名
	Signature string `json:"signature"` //个性签名
	Vip       int    `json:"vip"`       //是否是大大大会员
	Level     int    `json:"level"`     //几级号
	HeadPath  string `json:"head_path"` //头像地址
	Position  string `json:"position"`  //在作品中的职位，用于联合投稿的情况
}

//视频参加的活动结构
type DetailedVideoActivity struct {
	Join int    `json:"join"` //是否参加，默认为0(不参加)，1为参加
	Name string `json:"name"` //活动名字
	Url  string `json:"url"`  //活动URL
}

//视频的更多选项
type DetailedVideoMore struct {
	PersonalDeclaration  int `json:"p_declar"`   //自制声明，默认为0(不添加)，1为添加
	WaterMark            int `json:"water_mark"` //水印类型
	BusincessDeclaration int `json:"b_declar"`   //商业声明，默认为0(不含)，1为含
	SubtitleLanguage     int `json:"sub_lang"`   //字幕语言类型
	SubtitleOpen         int `json:"sub_open"`   //是否允许粉丝投稿字幕，默认为0(不允许)，1为允许
}

//视频信息结构
type VideoInformation struct {
	gorm.Model
	BvCode    string `json:"bv_code"`    //bv号
	CoverPath string `json:"cover_path"` //封面文件地址
	Title     string `json:"title"`      //标题
	DateTime  string `json:"date_time"`  //发布时间
	Brief     string `json:"brief"`      //简介
	Plays     int64  `json:"plays"`      //播放量
	Type      int    `json:"type"`       //视频的类型(自制还是转载)
	P         int    `json:"p"`          //该视频的分集数
}

//视频路径信息结构
type VideoPathInformation struct {
	gorm.Model
	P    int    `json:"p" gorm:"column:p"`             //第几分p
	Path string `json:"path" gorm:"column:video_path"` //视频地址
	Name string `json:"name" gorm:"column:video_name"` //视频名
}

//弹幕结构
type VideoBarrage struct {
	gorm.Model
	DateTime  string `json:"date_time"`  //弹幕发表日期
	VideoTime string `json:"video_time"` //弹幕在视频中出现的时间点
	UsersId   int64  `json:"users_id"`
	Content   string `json:"content"` //弹幕内容
	Type      int    `json:"type"`    //弹幕类型
	Size      int    `json:"size"`    //弹幕字体大小
	Pattern   int    `json:"pattern"` //弹幕飘出表现形式
	Color     int    `json:"color"`   //弹幕颜色
}

func (v VideoAPI) TableName() string {
	return "videos_information"
}

func (m MetaComment) TableName() string {
	return "videos_meta_comments"
}

func (r ReplyComment) TableName() string {
	return "videos_reply_comments"
}

func (o OperateVideoInformation) TableName() string {
	return "users_operate_videos_relationship"
}

func (v VideoAuthorInformation) TableName() string {
	return "users_information"
}

func (v VideoInformation) TableName() string {
	return "videos_information"
}

func (v VideoBarrage) TableName() string {
	return "videos_barrages"
}

func (v VideoPathInformation) TableName() string {
	return "videos_path"
}
