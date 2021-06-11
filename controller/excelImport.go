package controller

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"kxtGo/model"
	"kxtGo/services"
	"kxtGo/tool"
	"os"
	"strconv"
	"strings"
	"time"
)

type ExcelImportTrademarkController struct {

}


func (excel *ExcelImportTrademarkController) Router(Group *gin.RouterGroup)  {
		// 批量导入（商标转让）
		Group.POST("/excelImport", excel.GinHandler)
}

func (excel *ExcelImportTrademarkController)  GinHandler(c *gin.Context){
	userId,ok := c.GetPostForm("userId")
	if  !ok {
		tool.JsonError(c,"请传递userId")
		return
	}

	//接收文件
	file, _ := c.FormFile("file")

	osFile,err := file.Open()
	if err != nil {
		tool.JsonError(c,"打开文件失败:"+err.Error())
		return
	}

	excelFile, err := excelize.OpenReader(osFile)
	if err != nil {
		tool.JsonError(c,"打开表格失败:"+err.Error())
		return
	}

	excelFile.Path = file.Filename

	MysqlDb ,err := new(tool.Db).GetDb()
	if err != nil {
		tool.JsonError(c,"数据库连接失败:"+err.Error())
		return
	}

	// Get all the rows in the Sheet1.
	rows, err := excelFile.GetRows("Sheet1")
	var batchSlice []model.XsBrandSell
	timeUnix := time.Now().Unix() //秒级时间戳
	for rowK, row := range rows {
		if rowK>0 {
			XsBrandSellType := new(model.XsBrandSell)
		rowFor: for k, colCell := range row {
			if k <= 11 {
				colCell = strings.Replace(colCell, " ", "", -1)
				switch k{
				case 1,2,3,4,5,6,7,8,9:
					if colCell == "" {
						axis,_ := excelize.CoordinatesToCellName(13,rowK+1)
						_ = excelFile.SetCellValue("Sheet1", axis, "第"+strconv.Itoa(k)+"列参数不能为空")
						break rowFor
					}
				}

				switch k {
				case 1:
					XsBrandSellType.Title = colCell
				case 2:
					XsBrandSellType.Company = colCell
				case 3:
					XsBrandSellType.Country = colCell
				case 4:
					classIdStr := ""
					classIdStrDb := ""
					colCellList := strings.Split(colCell,",")
					for _,colCellListVal:= range colCellList {
						pos := strings.Index(colCellListVal,"-")
						if classIdStr != "" {
							classIdStr += ". "+colCellListVal[0:pos]
							classIdStrDb += ","+colCellListVal[0:pos]
						}else{
							classIdStr = colCellListVal[0:pos]
							infoModel := new(model.InfoModel)
							MysqlDb.Table("xs_info").Where("title",colCellListVal).Select("id").Limit(1).Find(infoModel)
							if infoModel.Id == 0 {
								axis,_ := excelize.CoordinatesToCellName(13,rowK+1)
								_ = excelFile.SetCellValue("Sheet1", axis, "第"+strconv.Itoa(k)+"列类目系统不存在")
								break rowFor
							}
							classIdStrDb = strconv.Itoa(infoModel.Id)
						}
					}
					XsBrandSellType.Types = classIdStrDb

					imagesMergeFunc := services.NewImagesTypeDefaultConfig()
					titleStr := XsBrandSellType.Country+"商标·第"+classIdStr+"类"
					classStr := XsBrandSellType.Title

					if classStr == "中国" {
						imagesMergeFunc.SetBgImage("zg.jpg")
					}else{
						imagesMergeFunc.SetBgImage("om.jpg")
					}
					mark:=row[10]
					if mark != "" {
						imagesMergeFunc.SetMark(mark+"标")
					}
					ser , err := imagesMergeFunc.ImageMerge(titleStr,classStr)
					if err != nil {
						axis,_ := excelize.CoordinatesToCellName(13,rowK+1)
						_ = excelFile.SetCellValue("Sheet1", axis, "第"+strconv.Itoa(k)+"列合并图片失败")
						break rowFor
					}

					XsBrandSellType.Logo = ser
				case 5:
					XsBrandSellType.No = colCell
				case 6:
					colCellFloat64,_ := strconv.ParseFloat(colCell,32)
					XsBrandSellType.Price = float32(colCellFloat64)
				case 7:
					XsBrandSellType.BrandSmall = colCell
				case 8:
					XsBrandSellType.ContactName = colCell
				case 9:
					XsBrandSellType.ContactPhone = colCell
				case 10:
					if colCell != "" {
						XsBrandSellType.Mark = colCell+"标"
					}
				case 11:
					XsBrandSellType.Cost = colCell
				}
			}
		}
			XsBrandSellType.UserId,_ = strconv.Atoi(userId)
			XsBrandSellType.CreateTime = timeUnix
			XsBrandSellType.UpdateTime = timeUnix

			batchSlice = append(batchSlice,*XsBrandSellType)
		}

		// 数据多了先存一波
		if len(batchSlice) > 50 {
			MysqlDb.Table("xs_brand_sell").Create(batchSlice)
			batchSlice = batchSlice[0:0] //清空
		}
	}

	if len(batchSlice) >0 {
		MysqlDb.Table("xs_brand_sell").Create(batchSlice)
	}

	outFileName := "import_" + strconv.FormatInt(timeUnix, 10) + file.Filename

	err = excelFile.SaveAs(outFileName)
	if err != nil {
		tool.JsonError(c,"保存反馈文件失败:"+err.Error())
		return
	}

	//结束删除文件
	defer os.Remove(outFileName)

	cosPath := "/temp/"+time.Now().Format("2006_01_02")+"/"+outFileName
	ossUrl,err := new(tool.Oss).UploadFile(cosPath,outFileName)
	if err != nil {
		tool.JsonError(c,"OSS上传失败:"+err.Error())
		return
	}

	tool.JsonSuccess(c,map[string]string{
		"url":ossUrl,
	})
}