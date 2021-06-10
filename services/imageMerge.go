package services

import (
	"errors"
	"github.com/fogleman/gg"
	"image"
	"kxtGo/tool"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type ImagesType struct {
	fontName string
	bgImg string
	mark string
}


func NewImagesTypeDefaultConfig() ImagesType {
	return ImagesType{
		fontName : "resource/font/simhei.ttf",
		bgImg : "resource/bgImages/om.jpg",
	}
}

//设置字体名称
func (i *ImagesType) SetFontName(fontName string) {
	i.fontName = "resource/font/"+fontName
}

//设置图片名称
func (i *ImagesType) SetBgImage(bgImg string) {
	i.bgImg = "resource/bgImages/"+bgImg
}

func (i *ImagesType) SetMark(mark string) {
	i.mark = mark
}


//合并图片
func (i *ImagesType) ImageMerge(titleStr string,classStr string) (string , error) {
	//打开背景图
	imageFile ,err := os.Open(i.bgImg)

	if err != nil {
		return "",errors.New("背景图打开失败:"+err.Error())
	}

	defer imageFile.Close()

	imageObj ,_,err := image.Decode(imageFile)

	if err != nil {
		return "",errors.New("图片解码失败:"+err.Error())
	}

	//加载背景图
	dc := gg.NewContextForImage(imageObj)

	//加载字体显示大小
	if err := dc.LoadFontFace(i.fontName, 40); err != nil {
		return "",errors.New("字体文件错误:"+err.Error())
	}

	//白色字体颜色
	dc.SetRGB(1, 1, 1)

	var str string
	var sWidth float64

	//第一行标题字体合成
	str = titleStr
	sWidth, _ = dc.MeasureString(str)
	dc.DrawString(str,(float64(dc.Width())-sWidth)/2 , 200)

	//第二行类目字体合成
	_ = dc.LoadFontFace(i.fontName, 60)
	str = classStr
	sWidth, _ = dc.MeasureString(str)
	dc.DrawString(str, (float64(dc.Width())-sWidth)/2 , 300)

	//第三行状态
	str = i.mark
	if str != "" {
		_ = dc.LoadFontFace(i.fontName, 27)
		sWidth, _ = dc.MeasureString(str)
		dc.DrawString(str, (float64(dc.Width())-sWidth)/2 , 375)
	}

	//保存图片
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	timeUnix := time.Now().UnixNano()
	outFileName := "merge_" + strconv.FormatInt(timeUnix, 10)+"_"+ strconv.Itoa(r.Intn(10)) + ".png"
	err = dc.SavePNG(outFileName)
	if err != nil {
		return "",errors.New("图片保存本地失败:"+err.Error())
	}

	//结束删除文件
	defer os.Remove(outFileName)

	cosPath := "/micro/merge/"+outFileName
	ossUrl,err := new(tool.Oss).UploadFile(cosPath,outFileName)
	if err != nil {
		return "",errors.New("OSS上传失败:"+err.Error())
	}

	return  ossUrl,nil
}