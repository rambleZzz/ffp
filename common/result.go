package common

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"time"
)

func ResultExport() {
	// 结果输出
	f := excelize.NewFile()
	ExcelFile = GetOutputExcelname()
	// 输出目标详情至终端 和excel
	fmt.Printf("\n【目标详情】Tips:因字段过多输出终端显示比较乱，在终端只显示部分内容，更多详情请在excel结果文件中查看(前提没有使用-noe参数).\n")
	if len(DomainScanResultMin) > 0 {
		DomainScanResultMin2DInterface, _ := ConvertStructSliceTo2D(DomainScanResultMin)
		Out2terminal(Header5, ConvertInterfaceSlice2DToStringSlice(DomainScanResultMin2DInterface))
		if !NotOutExcelFlag {
			OutDomainScanResult2Excel(DomainScanResult, f)
		}
	} else {
		OutNone2Terminal()
	}

	// 输出IP详情至终端 和 excel
	fmt.Printf("\n【无CDN IP统计】\n")
	if len(IPInfoResult) > 0 {
		ipInfoResult2DInterface, _ := ConvertStructSliceTo2D(IPInfoResult)
		Out2terminal(Header2, ConvertInterfaceSlice2DToStringSlice(ipInfoResult2DInterface))
		if !NotOutExcelFlag {
			OutipInfoResult2Excel(IPInfoResult, f)
		}
	} else {
		OutNone2Terminal()
	}

	// 输出无CDN C段IP地址至终端 和 excel
	fmt.Printf("\n【无CDN IP C段统计】\n")
	if len(IPResult) > 0 {
		cidrAddresses := GetCIDRAddresses(IPResult)
		Out2terminal(Header3, ConvertStringSliceTo2D(cidrAddresses))
		if !NotOutExcelFlag {
			OutcidrAddresses2Excel(cidrAddresses, f)
		}
	} else {
		OutNone2Terminal()
	}

	//输出扫描个数和未扫描个数至终端，未扫描目标至excel
	notAliveCount := len(NotAliveTargetResult)
	if notAliveCount > 0 && !NotOutExcelFlag {
		OutNotAliveTarget2Excel(NotAliveTargetResult, f)
	}
	fmt.Printf("[%v] All scans have been completed.\n", GetCurrentTimeStringTime())
	fmt.Printf("[%v] 共耗时[%v]扫描目标%v个，未扫描%v个. tips: 因未存活而未扫描的目标已输出至excel文档(前提没有使用-noe参数). \n", GetCurrentTimeStringTime(), time.Now().Sub(StartTime), len(AliveTargetResult), notAliveCount)
	if !NotOutExcelFlag {
		fmt.Printf("[%v] Result Path：%v\n", GetCurrentTimeStringTime(), ExcelFile)
	}
}

func OutNone2Terminal() {
	Out2terminal(Header0, [][]string{{"无"}})
}
