package common

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
)

// 输出目标结果到excel
func OutDomainScanResult2Excel(s interface{}, f *excelize.File) {
	data, _ := ConvertStructSliceTo2D(s)
	SaveExcel(SheetName1, Header1, data, f)
}

// 输出ip详情结果结果到excel
func OutipInfoResult2Excel(s interface{}, f *excelize.File) {
	data, _ := ConvertStructSliceTo2D(s)
	SaveExcel(SheetName2, Header2, data, f)
}

// 输出无CDN IP C段结果到excel
func OutcidrAddresses2Excel(s []string, f *excelize.File) {
	data := ConvertStringSliceTo2D(s)
	SaveExcel(SheetName3, Header3, ConvertStringSlice2DToInterface(data), f)
}

// 输出不存活目标结果到excel
func OutNotAliveTarget2Excel(s []string, f *excelize.File) {
	data := ConvertStringSliceTo2D(s)
	SaveExcel(SheetName4, Header4, ConvertStringSlice2DToInterface(data), f)
}

// 保存文件统一处理
func SaveExcel(sheetName string, header []string, data [][]interface{}, f *excelize.File) {
	f, _ = ExportExcel(sheetName, header, data, f)
	f.DeleteSheet("Sheet1")
	f.SaveAs(ExcelFile)
}

// 获取结果文件名
func GetOutputExcelname() string {
	cTime := GetCurrentTimeString()
	excelName := fmt.Sprintf("ffp_%s.xlsx", cTime)
	filename := filepath.Join(ResultPath, excelName)
	return filename
}

// 终端输出
func Out2terminal(keys []string, values [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	//table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader(keys)
	table.AppendBulk(values)
	table.Render()
}

// 创建一个进度条
func NewProgressbar(max int, description string) *progressbar.ProgressBar {
	d := fmt.Sprintf("[cyan][%s][reset] is running...", description)
	var bar = progressbar.NewOptions(max,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription(d),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}
