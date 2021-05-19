package model

import (
	"SecondTerm/Homework-7/videosRouters/dao"
	ut "SecondTerm/Homework-7/videosRouters/utilities"
	"gorm.io/gorm"
	"strings"
	"time"
)

//获取视频ID
func videoId(api *ut.VideoAPI) {
	dao.DB.Select("id").Where(api, "bv_code").Find(api)
}

//更改时间格式
func datetimeChange(s *string) {
	if *s == "" {
		return
	}

	*s = strings.Replace(*s, "T", " ", 1)
	*s = (*s)[0:19]
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

	err := dao.DB.Model(&v).Select("bv_code", "cover_path", "title", "date_time", "plays", "author_id").
		Where("match(title) against (? in natural language mode )", target).Limit(count).Find(&v).Error
	if err != nil {
		return nil, err
	}

	for i := range v {
		err = dao.DB.Model(&ut.VideoAuthorInformation{}).
			Select("name").
			Where("id = ?", v[i].AuthorId).
			Find(&v[i].Author).Error
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

//获取视频详情。
func GetDetailedVideoInformation(api *ut.VideoAPI) (ut.VideoInformation, error) {
	var (
		v ut.VideoInformation
	)
	videoId(api)

	//获取视频信息
	err := dao.DB.Model(&v).
		Omit("id", "created_at", "updated_at", "deleted_at").
		Where("id = ?", api.ID).
		Find(&v).Error
	if err != nil {
		return ut.VideoInformation{}, err
	}

	//获取作者信息
	err = dao.DB.Model(&ut.VideoAuthorInformation{}).
		Select("name", "signature", "vip", "level", "head_path").
		Where("id = ?", v.AuthorId).
		Find(&v.Author).Error
	if err != nil {
		return ut.VideoInformation{}, err
	}

	//获取点赞，硬币，收藏，分享数
	err = dao.DB.Table("users_operate_videos_relationship o").
		Select("sum(o.likes) as likes,sum(o.coins) as coins,sum(o.collections) as collections,sum(o.shares) as shares").
		Find(&v.Common).Error
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
		Where("videos_id = ? and p = ?", api.ID, api.P).
		Find(&vb).Error
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
		Where("videos_id = ?", api.ID).
		Find(&vpi).Error
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

	//获取置顶评论
	err := dao.DB.Model(&vcTop).Select("id", "date_time", "content", "floor", "author_id").
		Where("video_id = ? and top = 1", api.ID).
		Find(&vcTop).Error
	if err != nil {
		return ut.MetaComment{}, nil, nil, err
	}
	//获取其他信息
	err = commentOther(&vcTop, 0)
	if err != nil {
		return ut.MetaComment{}, nil, nil, err
	}

	//获取非置顶评论
	err = dao.DB.Model(&vcCommon).Order("id desc").
		Select("id", "date_time", "content", "floor", "author_id").
		Where("video_id = ? and top = 0", api.ID).
		Find(&vcCommon).Error
	if err != nil {
		return ut.MetaComment{}, nil, nil, err
	}
	//获取其他信息
	for i := range vcCommon {
		err := commentOther(&vcCommon[i], 1)
		if err != nil {
			return ut.MetaComment{}, nil, nil, err
		}
	}

	//分离热门评论
	vcCommon, vcHot = separateComment(vcCommon)

	return vcTop, vcHot, vcCommon, nil
}

//获取评论的其他信息
func commentOther(c *ut.MetaComment, commentType int) error {
	err := replyComment(c)
	if err != nil {
		return err
	}

	err = commentAuthor(c)
	if err != nil {
		return err
	}

	err = commentHeat(c, commentType)
	if err != nil {
		return err
	}

	return nil
}

//获取元评论的回复评论
func replyComment(c *ut.MetaComment) error {
	err := dao.DB.Select("id", "date_time", "content", "author_id", "reply_author_id").
		Where("reply_comment_id = ?", (*c).ID).
		Find(&(*c).ReplyComment).Error
	return err
}

//获取元评论及其回复评论的作者信息
func commentAuthor(c *ut.MetaComment) error {
	err := dao.DB.Select("id", "name", "vip", "level").
		Where("id = ?", (*c).AuthorId).
		Find(&(*c).Author).Error
	if err != nil {
		ut.LogError("GetcommentOther Error", err)
		return err
	}
	datetimeChange(&(*c).DateTime)

	//获取回复评论作者信息
	for i := range (*c).ReplyComment {
		dao.DB.Select("id", "name", "vip", "level").
			Where("id = ?", (*c).ReplyComment[i].AuthorId).
			Find(&(*c).ReplyComment[i].Author)

		dao.DB.Select("id", "name", "vip", "level").
			Where("id = ?", (*c).ReplyComment[i].ReplyAuthorId).
			Find(&(*c).ReplyComment[i].ReplyAuthor)

		datetimeChange(&(*c).ReplyComment[i].DateTime)
	}

	return nil
}

//获取评论热度
func commentHeat(c *ut.MetaComment, commentType int) error {
	err := dao.DB.Model(&gorm.Model{}).
		Table("likes_videos_comments_relationship").
		Select("sum(if(likes = 1, 1, 0)) as likes,sum(if(likes = -1, 1, 0)) as heat").
		Where("comments_id = ? and comments_type = ?", (*c).ID, commentType).
		Find(c).Error
	if err != nil {
		return err
	}

	(*c).Heat = (*c).Likes - (*c).Heat
	return nil
}

//分离热门评论
func separateComment(c []ut.MetaComment) ([]ut.MetaComment, []ut.MetaComment) {
	var (
		common []ut.MetaComment
		hot    []ut.MetaComment
	)

	for _, v := range c {
		//以10为分界线进行分类
		if v.Heat >= 10 {
			hot = append(hot, v)
		} else {
			common = append(common, v)
		}

		//综合点赞数升序整理
		sortCommentsByDate(&hot)
	}

	return common, hot
}

//整理非置顶评论。
//采用冒泡排序达到热度升序。
func sortCommentsByDate(data *[]ut.MetaComment) {
	var check int //检测是否进行交换，用于提前退出循环
	dataLen := len(*data)
	for i := 0; i < dataLen; i++ {
		for j := i + 1; j < dataLen; j++ {
			if (*data)[i].Heat < (*data)[j].Heat {
				(*data)[i], (*data)[j] = (*data)[j], (*data)[i]
				check++
			}
		}
		if check == 0 { //若在一次循环中没有交换，则已有序，退出循环
			break
		}
		check = 0 //重置判断变量
	}
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
		err := dao.DB.Table("videos_reply_comments").Omit("floor").Create(&nc).Error
		if err != nil {
			return err
		}
	}
	return nil
}
