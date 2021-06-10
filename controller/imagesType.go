package controller

import (
	"github.com/gin-gonic/gin"
	"kxtGo/services"
	"kxtGo/tool"
)

type ImagesTypeTrademarkController struct {
}

func (I *ImagesTypeTrademarkController) Router(Group *gin.RouterGroup)  {
	// 批量导入（商标转让）
	Group.POST("/imageMerge", I.GinHandler)
}

func  (I *ImagesTypeTrademarkController)  GinHandler(c *gin.Context)  {
	imagesMergeFunc := services.NewImagesTypeDefaultConfig()

	if fontName,ok := c.GetPostForm("fontName"); ok {
		imagesMergeFunc.SetFontName(fontName)
	}

	if bgImg ,ok := c.GetPostForm("bgImg"); ok {
		imagesMergeFunc.SetBgImage(bgImg)
	}

	if mark ,ok := c.GetPostForm("mark"); ok {
		imagesMergeFunc.SetMark(mark)
	}

	titleStr,ok1 := c.GetPostForm("titleStr")
	classStr,ok2 := c.GetPostForm("classStr")
	if !ok1 || !ok2 {
		tool.JsonError(c,"titleStr或classStr不存在")
		return
	}

	ser , err := imagesMergeFunc.ImageMerge(titleStr,classStr)
	if err != nil {
		tool.JsonError(c,err.Error())
		return
	}

	tool.JsonSuccess(c,ser)
	return
}