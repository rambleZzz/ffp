package domain_scan

import (
	"fmt"
	"github.com/rambleZzz/ffp/common"
	"strings"
	"sync"
)

type DomainScan struct {
	Domain           string             `json:"domain"`
	HttpxInfo        HttpxResult        `json:"httpxInfo"`
	CDNCheckInfo     CDNCheckResult     `json:"cdnCheckInfo"`
	ObserverWardInfo ObserverWardResult `json:"observerWardInfo"`
}

func NewDomainScan() *DomainScan {
	return &DomainScan{}
}

func (d *DomainScan) RunDomainScan(target string) {
	h := NewHttpx()
	c := NewCDNcheck()
	o := NewObserverWard()
	// httpx
	h.RunHttpx(target)
	// 先判断是否有httpx结果
	if !IsEmptyStruct(&h.HttpxResult) {
		prefixTarget := strings.Split(target, "/")[0]
		if len(h.HttpxResult.A) == 0 {
			h.HttpxResult.A = []string{h.HttpxResult.Host}
			// 判断目标前缀不为 ip 或 ipport格式
		} else if !common.ChcekIsIp(prefixTarget) && !common.ChcekIsIpPort(prefixTarget) {
			// cdnCheck CDN检测
			c.RunCDNcheck(target)
		}
		//  如果存活再ObserverWard指纹识别
		if h.HttpxResult.StatusCode != 0 {
			o.RunObserverWard(h.HttpxResult.Url)
			// iconHash 检测
			iconUrl := fmt.Sprintf("%v://%v", h.HttpxResult.Scheme, prefixTarget)
			iconHash := GetFaviconMmh3Hash(iconUrl)
			h.HttpxResult.FaviconHash = iconHash
		}

	}
	d.Domain = target
	d.HttpxInfo = h.HttpxResult
	d.CDNCheckInfo = c.CDNCheckResult
	d.ObserverWardInfo = o.ObserverWardResult
}

func (d *DomainScan) GetDomainDetail() (domainDetail common.DomainDetail) {
	domainDetail = common.DomainDetail{
		Domain:             d.Domain,
		Title:              d.HttpxInfo.Title,
		A:                  strings.Join(d.HttpxInfo.A, "|"),
		WebServer:          d.HttpxInfo.WebServer,
		ObserverWardFinger: strings.Join(d.ObserverWardInfo.Name, "|"),
		Tech:               strings.Join(d.HttpxInfo.Tech, "|"),
		Url:                d.HttpxInfo.Url,
		IsCDN:              d.CDNCheckInfo.IsCDN,
		CName:              strings.Join(d.CDNCheckInfo.Cname, "|"),
		CDNName:            strings.Join(d.CDNCheckInfo.CDNName, "|"),
		Host:               d.HttpxInfo.Host,
		Port:               d.HttpxInfo.Port,
		ContentType:        d.HttpxInfo.ContentType,
		StatusCode:         d.HttpxInfo.StatusCode,
		FinalUrl:           d.HttpxInfo.FinalUrl,
		Alive:              !d.HttpxInfo.Failed,
		Scheme:             d.HttpxInfo.Scheme,
		FaviconHash:        d.HttpxInfo.FaviconHash,
	}
	// 标题为空或为乱码重新赋值
	if (domainDetail.Title == "" && d.ObserverWardInfo.Title != "") || common.HasInvalidChars(domainDetail.Title) {
		domainDetail.Title = d.ObserverWardInfo.Title
	}
	// A记录主机大于3 默认为有CDN
	if domainDetail.StatusCode != 0 && domainDetail.IsCDN == false {
		if len(d.HttpxInfo.A) > 3 {
			domainDetail.IsCDN = true
			domainDetail.CDNName = "CDN_IP"
		} else {
			// 统计所有不是CDN的IP
			common.IPResult = append(common.IPResult, d.HttpxInfo.A...)
		}
	}
	return domainDetail
}

func (d *DomainScan) Do(domain string, wg *sync.WaitGroup) {
	var domainDetail common.DomainDetail
	d.RunDomainScan(domain)
	domainDetail = d.GetDomainDetail()
	//fmt.Println(domainDetail)
	common.DomainScanResult = append(common.DomainScanResult, domainDetail)
	common.DomainScanResultMin = append(common.DomainScanResultMin, common.DomainDetailMin{
		Url:                domainDetail.Url,
		Title:              domainDetail.Title,
		WebServer:          domainDetail.WebServer,
		ObserverWardFinger: domainDetail.ObserverWardFinger,
		Tech:               domainDetail.Tech,
		IsCDN:              domainDetail.IsCDN,
		StatusCode:         domainDetail.StatusCode,
	})
}

// 判断结构体的值是否为空
func IsEmptyStruct(s *HttpxResult) bool {
	return s == nil
}
