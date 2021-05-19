package serve

import (
	"SecondTerm/Homework-7/videosRouters/model"
	ut "SecondTerm/Homework-7/videosRouters/utilities"
	"github.com/gin-gonic/gin"
	"strconv"
)

//获取视频评论。
func GetVideoComments(ctx *gin.Context) {
	var (
		resp           ut.Resp
		topComment     ut.MetaComment   //置顶评论
		hotComments    []ut.MetaComment //热门评论
		commonComments []ut.MetaComment //平平无奇的评论
		api            ut.VideoAPI
	)

	var ()

	api.BVCode = ctx.Query("bv_code")
	topComment, hotComments, commonComments, err := model.VideoComments(&api) //获取非置顶评论
	if err != nil {
		resp.Return(ctx, 50001, "未知错误", nil)
		return
	}

	resp.ReturnMuti(ctx, 500, "响应成功", "top_comment", topComment, "hot_comment", hotComments,
		"common_comment", commonComments)
}

//获取视频信息。
func GetVideoInformation(ctx *gin.Context) {
	var resp ut.Resp

	videoType := ctx.Query("type")

	if videoType == "1" { //视频的简要信息
		target := ctx.Query("target")
		count, _ := strconv.Atoi(ctx.Query("count"))

		data, err := model.GetBriefVideoInformation(target, count)
		if err == nil {
			resp.Return(ctx, 600, "响应成功", data)
		} else {
			resp.Return(ctx, 60002, "未知错误", nil)
		}
	} else if videoType == "2" { //视频的详细信息
		var api ut.VideoAPI
		api.BVCode = ctx.Query("bv_code")

		data, err := model.GetDetailedVideoInformation(&api)
		if err == nil {
			resp.ReturnAll(ctx, 600, "响应成功", data)
		} else {
			resp.Return(ctx, 60002, "未知错误", nil)
		}
	} else { //视频的类型无法识别
		resp.Return(ctx, 60001, "视频信息的类型不合法", nil)
	}
}

//获取视频弹幕。
func GetVideoBarrages(ctx *gin.Context) {
	var (
		resp ut.Resp
		api  ut.VideoAPI
	)

	api.BVCode = ctx.Query("bv_code")
	api.P, _ = strconv.Atoi(ctx.Query("p"))

	data, err := model.GetVideoBarrages(&api)
	if err == nil {
		resp.ReturnAll(ctx, 700, "响应成功", data)
	} else {
		resp.Return(ctx, 70001, "未知错误", nil)
	}
}

//获取视频地址。
func GetVideoPath(ctx *gin.Context) {
	var (
		resp ut.Resp
		api  ut.VideoAPI
	)

	api.BVCode = ctx.Query("bv_code")

	data, err := model.GetVideoPath(&api)

	if err == nil {
		resp.Return(ctx, 800, "响应成功", data)
	} else {
		resp.Return(ctx, 80001, "响应成功", nil)
		resp.Message = err.Error()
	}
}

//进行视频操作。
func OperateVideo(ctx *gin.Context) {
	var (
		resp ut.Resp
		op   ut.OperateVideoInformation
		err  error
	)

	op.VideoId, _ = strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	op.UserId, _ = strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	operateType := ctx.Query("type")             //进行了什么操作
	value, _ := strconv.Atoi(ctx.Query("value")) //操作后的数值一

	if op.VideoId == 0 || op.UserId == 0 {
		resp.Return(ctx, 90001, "信息不完整", nil)
		return
	}

	if operateType == "1" { //点赞操作
		if value > 1 || value < -1 { //如果点赞数值不合法
			resp.Return(ctx, 90002, "点赞数值不合法", nil)
			return
		}
		op.Like = value

	} else if operateType == "2" { //投币操作
		if value > 2 || value < 1 { //如果投币数值不合法
			resp.Return(ctx, 90003, "投币数值不合法", nil)
			return
		}
		op.Coin = value

	} else if operateType == "3" { //收藏操作
		if value > 1 || value < 0 { //如果收藏数值不合法
			resp.Return(ctx, 90004, "收藏数值不合法", nil)
			return
		}
		op.Collect = value

	} else if operateType == "4" { //分享操作
		if value > 1 || value < 0 { //如果分享数值不合法
			resp.Return(ctx, 90005, "分享数值不合法", nil)
			return
		}
		op.Share = value

	} else if operateType == "6" { //一键三连操作
		op.Like = 1
		op.Coin = 1
		op.Collect = 1

	} else { //操作类型不合法
		resp.Return(ctx, 90002, "操作类型不合法", nil)
		return
	}

	err = model.UpdateVideoOperation(&op)
	if err == nil {
		resp.Return(ctx, 900, "响应成功", nil)
	} else {
		resp.Return(ctx, 90006, "未知错误", nil)
	}
}

//添加评论
func AddComment(ctx *gin.Context) {
	var (
		nc   ut.NewComment
		resp ut.Resp
	)
	err := ctx.ShouldBindJSON(&nc)
	if err != nil {
		ut.LogError("GetNewComment Error", err)
		resp.Return(ctx, 13001, "未知错误", nil)
		return
	}

	err = model.AddNewComment(&nc)
	if err != nil {
		ut.LogError("AddNewComment Error", err)
		resp.Return(ctx, 13001, "未知错误", nil)
	} else {
		resp.Return(ctx, 1300, "评论成功", nil)
	}
}
