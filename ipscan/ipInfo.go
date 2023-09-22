package ipscan

import (
	"github.com/rambleZzz/ffp/common"
	"log"
)

type IPInfo struct {
	IP      string
	Country string
	City    string
}

// NewIPInfo 创建NewIPInfo对象
func NewIPInfo() *IPInfo {
	return &IPInfo{}
}

func (i *IPInfo) RunIPInfo(ip string) {
	q, _ := NewQQwry(common.QqwryPath)
	result, err := q.Find(ip)
	if err != nil {
		log.Fatal(err)
	}
	common.IPInfoResult = append(common.IPInfoResult, common.IPInfo{
		IP:      ip,
		Country: result.Country,
		City:    result.City,
	})
}
