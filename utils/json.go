package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"time"
)

func OutJsonOk(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"rtn": 0,
		"msg": msg,
	})
}

func OutJson(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"rtn":  0,
		"msg":  "成功",
		"data": data,
	})
}

func OutParamErrorJson(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": 1,
		"msg": "参数错误",
	})
}

func OutErrorJson(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"rtn": 1,
		"msg": err.Error(),
	})
}

func OutErrorCodeJson(c *gin.Context, code int64, err error) {
	c.JSON(200, gin.H{
		"rtn": code,
		"msg": err.Error(),
	})
}

// OutAuthNeedError 需要登录
func OutAuthNeedError(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": 4001,
		"msg": "该功能需要登录!",
	})
}

// OutAuthOutdatedError 输出过期错误
func OutAuthOutdatedError(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": 4002,
		"msg": "登录态已过期,请重新登录!",
	})
}

// OutErrorJsonWithStr 自定义文本错误
func OutErrorJsonWithStr(c *gin.Context, err string) {
	c.JSON(200, gin.H{
		"rtn": 1,
		"msg": err,
	})
}

// OutRBACError 输出权限错误
func OutRBACError(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": -1,
		"msg": "该功能需要权限，请申请权限",
	})
}

// ValidateJson 监测json数据正确性
func ValidateJson(jsonString string) (right bool, data interface{}) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonString)))
	decoder.UseNumber()
	err := decoder.Decode(&data)
	if err != nil {
		return false, data
	}
	return true, data
}

// ExportExcel 导出Excel文件
func ExportExcel(c *gin.Context, file *excelize.File) {
	if file == nil {
		OutErrorJsonWithStr(c, "没有可导出的数据")
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"Workbook.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	_ = file.Write(c.Writer)
}

func ExportExcelTitle() string {
	timeStr := time.Now().Format("2006-01-02 15:04")
	return "数据详情（导出时间：" + timeStr + ")"
}

// ExportTitleStyle 标题样式
func ExportTitleStyle(f *excelize.File) int {
	style, err := f.NewStyle(`{
    "font":
    	{
       	 	"bold": true,
        	"family": "Times New Roman",
        	"size": 30,
        	"color": "#232323"
    	},
	"alignment":
    	{
       	 	"horizontal": "center",
        	"vertical": "center"
    	 }
	}`)
	if err != nil {
		fmt.Println("title", err)
	}
	return style
}

func ExportHeaderStyle(f *excelize.File) int {
	style, err := f.NewStyle(`{
    "alignment":
    	{
        	"horizontal": "center",
        	"vertical": "center"
     	},
	"fill":
		{
			"type":  "pattern",
			"color": ["#C4C4C4"],
			"pattern": 1
		},
	"font":
    	{
       	 	"bold": true,
        	"family": "Times New Roman",
        	"size": 15,
        	"color": "#000000"
    	}
  	}`)
	if err != nil {
		fmt.Println("header", err)
	}
	return style
}

// ExportDefaultAlignmentStyle 默认样式
func ExportDefaultAlignmentStyle(f *excelize.File) int {
	style, err := f.NewStyle(`{
    "alignment":
    	{
        	"horizontal": "left",
        	"vertical": "center",
			"wrap_text": true
     	}
 	 }`)
	if err != nil {
		fmt.Println("default", err)
	}
	return style
}

// AddExcelStyle 给excel表格添加样式
func AddExcelStyle(f *excelize.File, sheet, fromCol, toCol string) (startRow int) {
	_ = f.SetColStyle(sheet, fromCol+":"+toCol, ExportDefaultAlignmentStyle(f))
	_ = f.MergeCell(sheet, fromCol+"1", toCol+"4")
	_ = f.SetCellStyle(sheet, fromCol+"1", toCol+"4", ExportTitleStyle(f))
	_ = f.SetCellStr(sheet, fromCol+"1", ExportExcelTitle())

	_ = f.SetCellStyle(sheet, fromCol+"5", toCol+"5", ExportHeaderStyle(f))
	_ = f.SetRowHeight(sheet, 5, 30)

	PrintExcel(f, sheet)
	return 5
}

func PrintExcel(f *excelize.File, sheetName string) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
