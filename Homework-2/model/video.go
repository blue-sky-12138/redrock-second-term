package model

import (
	"Homework-2/dao"
	ut "Homework-2/utilities"
	"strings"
	"time"
)

//获取视频ID
func videoId(api *ut.VideoAPI) {
	dao.DB.Select("id").Where(api, "bv_code").Find(api)
}

//更改时间格式
func datetimeChange(s string) string {
	if s == "" {
		return ""
	}

	s = strings.Replace(s, "T", " ", 1)
	return s[0:19]
}

//获取视频简要信息。
//count表示一共要获取几条数据，最多获取1000条，若为0则默认获取1000条。
//target为搜索关键词。
func GetBriefVideoInformation(target string, count int) ([]ut.VideoInformation, error) {
	var (
		v []ut.VideoInformation
	)

	if count >= 0 || count > 1000 {
		count = 1000
	}

	err := dao.DB.Model(&v).Select("bv_code", "cover_path", "title", "date_time", "plays").
		Where("match(title) against (? in natural language mode )", target).Limit(count).Find(&v).Error
	if err != nil {
		return nil, err
	}
	return v, nil
}

//获取视频详情。
func GetDetailedVideoInformation(api *ut.VideoAPI) (ut.VideoInformation, error) {
	var (
		v ut.VideoInformation
	)
	videoId(api)

	err := dao.DB.Model(&v).Omit("id", "created_at", "updated_at", "deleted_at").Find(&v).Error
	if err != nil {
		return ut.VideoInformation{}, err
	}

	return v, nil
}

//获取视频弹幕。
func GetVideoBarrages(api *ut.VideoAPI) ([]ut.VideoBarrage, error) {
	var (
		vb []ut.VideoBarrage
	)
	videoId(api)

	err := dao.DB.Omit("videos_id", "p", "created_at", "updated_at", "deleted_at").
		Where("videos_id = ?", api.ID).Where("p = ?", api.P).Find(&vb).Error
	if err != nil {
		return nil, err
	}
	return vb, nil
}

//获取视频地址及相关信息。
func GetVideoPath(api *ut.VideoAPI) ([]ut.VideoPathInformation, error) {
	var (
		vpi []ut.VideoPathInformation
	)
	videoId(api)

	err := dao.DB.Select("p", "video_path", "video_name").
		Where("videos_id = (?)", api.ID).Find(&vpi).Error
	if err != nil {
		return nil, err
	}
	return vpi, nil
}

//获取视频评论
func VideoComments(api *ut.VideoAPI) (ut.MetaComment, []ut.MetaComment, []ut.MetaComment, error) {
	var (
		vcTop    ut.MetaComment
		vcHot    []ut.MetaComment
		vcCommon []ut.MetaComment
	)
	videoId(api)

	dao.DB.Model(&vcTop).Select("id", "date_time", "content", "floor", "author_id").
		Where("video_id = ? and top = 1", api.ID).Find(&vcTop)
	vcTop = commentOther(vcTop)

	dao.DB.Model(&vcCommon).Order("id desc").Select("id", "date_time", "content", "floor", "author_id").
		Where("video_id = ? and top = 0", api.ID).Find(&vcCommon)
	l := len(vcCommon)
	for i := 0; i < l; i++ {
		vcCommon[0] = commentOther(vcCommon[0])
	}

	return vcTop, vcHot, vcCommon, nil
}
func commentOther(c ut.MetaComment) ut.MetaComment {
	var (
		r []ut.ReplyComment
	)
	dao.DB.Select("id", "name", "vip", "level").Where("id = ?", c.AuthorId).Find(&c.Author)
	c.DateTime = datetimeChange(c.DateTime)

	//搜索回复评论
	dao.DB.Order("id desc").Select("id", "date_time", "content", "author_id", "reply_author_id").
		Where("reply_comment_id = ?", c.ID).Find(&r)
	l := len(r)
	for i := 0; i < l; i++ {
		dao.DB.Select("id", "name", "vip", "level").Where("id = ?", r[i].AuthorId).Find(&r[i].Author)
		dao.DB.Select("id", "name", "vip", "level").Where("id = ?", r[i].ReplyAuthorId).Find(&r[i].ReplyAuthor)
		r[i].DateTime = datetimeChange(r[i].DateTime)

	}

	c.ReplyComment = r
	return c
}

//更新用户对视频的操作
func UpdateVideoOperation(op *ut.OperateVideoInformation) error {
	err := dao.DB.Create(&op).Error
	if err != nil {
		return err
	}
	return nil
}

//添加新评论
func AddNewComment(nc *ut.NewComment) error {
	//赋与时间
	nc.DateTime = time.Now().Format("2006-01-02 15:04:05")

	if nc.ReplyAuthorId == 0 && nc.ReplyCommentId == 0 {
		dao.DB.Table("videos_meta_comments").Select("floor").Where("video_id = ?", nc.VideoId).Last(&nc)
		nc.Floor++

		err := dao.DB.Table("videos_meta_comments").Omit("reply_author_id", "reply_comment_id").Create(&nc).Error
		if err != nil {
			return err
		}
	} else {
		err := dao.DB.Table("videos_reply_comments").Omit("floor").Updates(&nc).Error
		if err != nil {
			return err
		}
	}
	return nil
}
