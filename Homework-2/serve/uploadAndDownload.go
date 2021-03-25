package serve

import (
	"Homework-2/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
	"time"
)

import (
	ut "Homework-2/utilities"
	"path"
)

//获取视频文件及封面
func GetVideoFile(ctx *gin.Context) {
	var (
		resp ut.Resp
		tem  ut.GetVideoFile
	)
	ctx.ShouldBindUri(&tem)

	if tem.BvCode == "" || tem.FileName == "" {
		resp.Return(ctx, 401, "参数错误!", nil)
	} else {
		filePath := path.Join("./static", "videos", tem.BvCode, tem.FileName)
		ctx.File(filePath)
	}
}

//获取用户头像
func GetUserHead(ctx *gin.Context) {
	var (
		resp ut.Resp
		tem  ut.GetUserHead
	)
	ctx.ShouldBindUri(&tem)
	if tem.Id == "" || tem.FileName == "" {
		resp.Return(ctx, 401, "参数错误!", nil)
	} else {
		headPath := path.Join("./static", "users", tem.Id, tem.FileName)
		ctx.File(headPath)
	}
}

////更新用户头像
func UpdateUserHead(ctx *gin.Context) {
	var resp ut.Resp

	temId, _ := strconv.ParseUint(ctx.PostForm("user_id"), 10, 64)
	userId := uint(temId)
	if userId == 0 {
		resp.Return(ctx, 11001, "用户ID不合法", nil)
		return
	}

	file, header, err := ctx.Request.FormFile("img")
	if err != nil {
		ut.LogError("GetFormFile(head) Error", err)
		resp.Return(ctx, 11002, "获取文件失败", nil)
		return
	}

	//获取文件后缀名
	var (
		fileNameLen = len(header.Filename) //文件名总长
		suffix      string
	)
	for i := fileNameLen; ; i-- {
		if header.Filename[i-1:i] == "." {
			suffix = header.Filename[i-1 : fileNameLen]
			break
		}
	}

	//路径添加，以系统Unix时间作为文件名
	newPath := fmt.Sprintf("%d/%d", userId, time.Now().Unix()) + suffix
	headPath := ut.HeadPath + newPath

	//保存文件
	out, err := os.Create(headPath)
	if err != nil {
		ut.LogError("CreateFile(head) Error", err)
		resp.Return(ctx, 11003, "未知错误11003", nil)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		ut.LogError("CopyFile(head) Error", err)
		resp.Return(ctx, 11003, "未知错误11003", nil)
		return
	}

	//删除用户原头像
	previousPath, err := model.GetUserHeadPath(userId)
	if err != nil { //获取用户原头像路径失败
		ut.LogError("GetFile(head) Error", err)
		resp.Return(ctx, 11004, "未知错误11004", nil)
		return
	}
	err = os.Remove(ut.HeadPath + previousPath)
	if err != nil { //删除用户原头像路失败
		ut.LogError("DeleteFile(head) Error", err)
		resp.Return(ctx, 11005, "未知错误11005", nil)
		return
	}

	err = model.ChangeUserHead(userId, newPath)
	if err != nil { //更新路径失败
		ut.LogError("UpdateFile(head) Error", err)
		resp.Return(ctx, 11006, "未知错误11006", nil)
		return
	}

	resp.Return(ctx, 1100, "更新成功", nil)
}

////Attention!未完工状态
////上传单个视频
//func UploadVideoOne(context *gin.Context) {
//	var (
//		resp utilities.Resp
//		data utilities.NewVideo
//	)
//
//	fileHeader, err := context.FormFile("file")
//	if err != nil {
//		utilities.LogError("GetVideoFile Error", err)
//		resp.Code = 13002
//		resp.Message = "未知错误13002"
//		context.JSON(200, resp)
//		return
//	}
//
//	//获取文件后缀名
//	var (
//		fileNameLen = len(fileHeader.Filename) //文件名总长
//		suffix      string
//	)
//	for i := fileNameLen; ; i-- {
//		if fileHeader.Filename[i-1:i] == "." {
//			suffix = fileHeader.Filename[i-1 : fileNameLen]
//			break
//		}
//	}
//
//	data.BvCode = utilities.NewBvCode()
//	videoPath := utilities.VideoPath + fileHeader.Filename + suffix
//
//	err = context.SaveUploadedFile(fileHeader, videoPath)
//	if err != nil {
//		utilities.LogError("SaveVideoFileError", err)
//		resp.Code = 13003
//		resp.Message = "未知错误13003"
//		context.JSON(200, resp)
//		return
//	}
//
//	err = model.AddNewVideoOne(&data)
//	if err != nil {
//		utilities.LogError("AddNewVideoOneInformation Error", err)
//		resp.Code = 13004
//		resp.Message = "未知错误13004"
//		context.JSON(200, resp)
//		return
//	}
//
//	resp.Code = 1300
//	resp.Message = "投稿成功"
//	context.JSON(200, resp)
//}

//Attention!未完工状态
//上传多个视频
func UploadVideoMore(context *gin.Context) {
	//var resp utilities.Resp

}

//Attention!未完工状态
//仅在本包使用。
//保存文件的快捷方式。
//func saveFile(context *gin.Context) error {
//
//}
