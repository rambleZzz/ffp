package ipscan

import (
	"fmt"
	"github.com/rambleZzz/ffp/common"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type IPInfo struct {
	IP        string
	Address   string
	Operator  string
	DataTwo   string
	DataThree string
	URL       string
}

// NewIPInfo 创建NewIPInfo对象
func NewIPInfo() *IPInfo {
	return &IPInfo{}
}

func (i *IPInfo) RunIPInfo(ip string) {
	ipInfo := i.GetIPInfo(ip)
	common.IPInfoResult = append(common.IPInfoResult, common.IPInfo{
		IP:        ip,
		Address:   ipInfo.Address,
		Operator:  ipInfo.Operator,
		DataTwo:   ipInfo.DataTwo,
		DataThree: ipInfo.DataThree,
		URL:       ipInfo.URL,
	})
}

func (i *IPInfo) GetIPInfo(ip string) IPInfo {
	url := fmt.Sprintf("http://www.cip.cc/%s", ip)

	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	// 提取 <pre> 与 </pre> 之间的内容
	content := extractContent(string(body), "<pre>", "</pre>")
	if err != nil {
		log.Fatal(err)
	}

	// 解析数据到结构体
	ipInfo := i.parseIPInfo(content)
	if err != nil {
		log.Fatal(err)
	}
	// 打印转换后的结构体
	fmt.Printf("\n%v\n", ipInfo)
	return ipInfo
}

// 提取指定开始和结束标记之间的内容
func extractContent(input, start, end string) string {
	startIndex := strings.Index(input, start) + len(start)
	endIndex := strings.Index(input, end)
	if startIndex == -1 || endIndex == -1 {
		return ""
	}
	return input[startIndex:endIndex]
}

// 解析数据到struct
func (i *IPInfo) parseIPInfo(str string) IPInfo {
	var ipInfo IPInfo
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "IP":
				ipInfo.IP = value
			case "地址":
				ipInfo.Address = value
			case "运营商":
				ipInfo.Operator = value
			case "数据二":
				ipInfo.DataTwo = value
			case "数据三":
				ipInfo.DataThree = value
			case "URL":
				ipInfo.URL = value
			default:
				fmt.Println("未知的键:", key)
			}
		}
	}
	return ipInfo
}
