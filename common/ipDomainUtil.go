package common

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// C段格式整理（多IP）
func GetCIDRAddresses(ipAddresses []string) []string {
	cidrMap := make(map[string]bool)
	for _, ipAddress := range ipAddresses {
		_, ipNet, err := net.ParseCIDR(ipAddress + "/24")
		if err != nil {
			fmt.Println("无效的 IP 地址:", ipAddress, err)
			continue
		}
		cidrMap[ipNet.String()] = true
	}
	cidrAddresses := make([]string, 0, len(cidrMap))
	for cidr := range cidrMap {
		cidrAddresses = append(cidrAddresses, cidr)
	}
	return cidrAddresses
}

// C段格式整理（单IP)
func GetCIDRAddres(ipAddress string) (string, error) {
	_, ipNet, err := net.ParseCIDR(ipAddress + "/24")
	if err != nil {
		fmt.Println("无效的 IP 地址:", ipAddress, err)
	}
	return ipNet.String(), err
}

// CheckDomain domain正则
func CheckDomain(domain string) bool {
	domainPattern := `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`
	reg := regexp.MustCompile(strings.TrimSpace(domainPattern))
	return reg.MatchString(domain)
}

// CheckUrlFormat url格式检查
func CheckUrlFormat(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}

// HostStrip domain处理
func HostStrip(u string) string {
	p, err := url.Parse(u)
	if err == nil && p.Host != "" {
		return p.Hostname()
	}
	host := strings.ReplaceAll(u, "https://", "")
	host = strings.ReplaceAll(host, "http://", "")
	host = strings.ReplaceAll(host, "/", "")
	host = strings.Split(host, ":")[0]
	return host
}

// url转domain，如 http://www.baidu.com/a/d > www.baidu.com
func GetDomainFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	domain := parsedURL.Hostname()
	return domain, nil
}

// DomainStrip domain处理
func DomainStrip(s string) string {
	domain := strings.TrimSpace(s)
	domain = strings.ReplaceAll(domain, "https://", "")
	domain = strings.ReplaceAll(domain, "http://", "")
	domain = strings.Split(domain, "/")[0]
	return domain
}

// TargetStrip domain处理
func TargetStrip(s string) string {
	target := strings.TrimSpace(s)
	target = strings.ReplaceAll(target, "https://", "")
	target = strings.ReplaceAll(target, "http://", "")
	target = strings.TrimSuffix(target, "/")
	return target
}

// 判断是否为ip
func ChcekIsIp(ip string) bool {
	s := strings.Trim(ip, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, s); match {
		return true
	}
	return false
}

// 判断是否为ip:port
func ChcekIsIpPort(iport string) bool {
	s := strings.Trim(iport, " ")
	regStr := `^(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})(\.(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})){3}:{1}(\d+)(,\d*)*$`
	if match, _ := regexp.MatchString(regStr, s); match {
		return true
	}
	return false
}

// 判断是否为192.168.1.1/24 192.168.1.1/16
func ChcekIsIpmask(ipmask string) bool {
	s := strings.Trim(ipmask, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])/([1-2][0-9]|3[0-2]|[1-9])$`
	if match, _ := regexp.MatchString(regStr, s); match {
		return true
	}
	return false
}

func CheckIPV4(ip string) bool {
	ipReg := `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`
	r, _ := regexp.Compile(ipReg)
	return r.MatchString(ip)
}
