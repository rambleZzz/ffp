package common

import "time"

// domainscan
var (
	DomainScanResult     []DomainDetail
	DomainScanResultMin  []DomainDetailMin
	IPResult             []string
	IPInfoResult         []IPInfo
	AliveTargetResult    []string
	NotAliveTargetResult []string
)

// task
var (
	StartTime time.Time
)

// excel输出
var (
	SheetName1 = "目标详情"
	SheetName2 = "IP详情"
	SheetName3 = "C段IP统计"
	SheetName4 = "不存活目标统计"
	Header1    = []string{"目标", "标题", "应用服务器", "ObserverWard指纹", "wappalyzer指纹", "网址", "是否有CDN", "A记录", "CNAME值", "CDN名称", "主机", "端口", "类型", "状态码", "icon哈希", "是否存活", "网站协议", "跳转URL"}
	Header2    = []string{"IP地址", "国家地区", "运营商"}
	Header3    = []string{"C段IP"}
	Header4    = []string{"不存活目标"}
	Header5    = []string{"网址", "标题", "应用服务器", "ObserverWard指纹", "wappalyzer指纹", "是否有CDN", "状态码"}
	Header0    = []string{"结果"}
	ExcelFile  string
)

// cli
var (
	NotOutExcelFlag bool
	Threads         int
)

// 路径
var (
	ResultPath     string
	CurrentRunPath string
)

// 域名详情
type DomainDetail struct {
	Domain             string `json:"domain"`
	Title              string `json:"title"`
	WebServer          string `json:"webserver"`
	ObserverWardFinger string `json:"observerWardFinger"`
	Tech               string `json:"tech"`
	Url                string `json:"url"`
	IsCDN              bool   `json:"isCDN"`
	A                  string `json:"a"`
	CName              string `json:"cname"`
	CDNName            string `json:"CDNName"`
	Host               string `json:"host"`
	Port               string `json:"port"`
	ContentType        string `json:"content_type"`
	StatusCode         int    `json:"status_code"`
	FaviconHash        string `json:"favicon_hash"`
	Alive              bool   `json:"alive"`
	Scheme             string `json:"scheme"`
	FinalUrl           string `json:"final_url"`
}

// 域名详情简化（输出至终端使用）
type DomainDetailMin struct {
	Url                string `json:"url"`
	Title              string `json:"title"`
	WebServer          string `json:"webserver"`
	ObserverWardFinger string `json:"observerWardFinger"`
	Tech               string `json:"tech"`
	IsCDN              bool   `json:"isCDN"`
	StatusCode         int    `json:"status_code"`
}

// IP归属信息
type IPInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
}

// url存活
type UrlAlive struct {
	Target  string `json:"target"`
	Url     string `json:"url"`
	IsAlive bool   `json:"is_alive"`
}

const banner = `
    ________________ 
   / ____/ ____/ __ \
  / /_  / /_  / /_/ /
 / __/ / __/ / ____/ 
/_/   /_/   /_/       Ver:1.2
`
const ffpAbout = `https://github.com/rambleZzz/ffp
FFP v1.2   (FastFingerPrint)  Dev:go1.20

`

// thirdparth
var (
	ObserverWardPath   string
	ObserverWardDir    string
	WebFingerprintPath string
	GeoLite2Path       string
	QqwryPath          string
)

// thirdparth
//var (
//	ObserverWardPath   = "/thirdparty/observerWard/observer_ward"
//	ObserverWardDir    = "/thirdparty/observerWard/"
//	WebFingerprintPath = "/thirdparty/observerWard/web_fingerprint_v3.json"
//	GeoLite2Path       = "/thirdparty/cdnCheck/GeoLite2-ASN.mmdb"
//)
