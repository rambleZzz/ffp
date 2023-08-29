package domain_scan

import (
	"crypto/tls"
	"fmt"
	"github.com/rambleZzz/ffp/common"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)


type UrlAlive struct {
	Client  *http.Client `json:"client"`
	Target  string       `json:"target"`
	Url     string       `json:"url"`
	IsAlive bool         `json:"is_alive"`
}

func NewUrlAlive() *UrlAlive {
	return &UrlAlive{}
}

func (u *UrlAlive) Do(target string, wg *sync.WaitGroup) {
	aliveInfo := u.GetTargetAliveInfo(target)
	if aliveInfo.IsAlive {
		common.AliveTargetResult = append(common.AliveTargetResult, aliveInfo.Target)
	} else {
		common.NotAliveTargetResult = append(common.NotAliveTargetResult, aliveInfo.Target)
	}
}

func (u *UrlAlive) GetAllTargetAliveInfo(targets []string) ([]UrlAlive, []UrlAlive) {
	var aliveTargets []UrlAlive
	var notAliveTargets []UrlAlive
	for _, t := range targets {
		target := common.TargetStrip(t)
		url, ok := u.isUrlAlive(target)
		urlAlive := UrlAlive{
			Target:  target,
			Url:     url,
			IsAlive: ok,
		}
		if ok {
			aliveTargets = append(aliveTargets, urlAlive)
		} else {
			notAliveTargets = append(notAliveTargets, urlAlive)
		}
	}
	return aliveTargets, notAliveTargets
}

func (u *UrlAlive) GetTargetAliveInfo(target string) UrlAlive {
	target = common.TargetStrip(target)
	url, ok := u.isUrlAlive(target)
	tAlive := UrlAlive{
		Target:  target,
		Url:     url,
		IsAlive: ok,
	}
	return tAlive
}

func (u *UrlAlive) GetHttpClient() *http.Client {
	to := 10000
	timeout := time.Duration(to * 1000000) //10秒
	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       timeout,
	}
	return client
}

func (u *UrlAlive) isListening(client *http.Client, url, method string) bool {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false
	}
	req.Header.Add("Connection", "close")
	req.Close = true
	resp, err := client.Do(req)
	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	if err != nil {
		return false
	}
	return true
}

func (u *UrlAlive) isUrlAlive(target string) (string, bool) {
	u.Client = u.GetHttpClient()
	url, ok := u.isHttpsAlive(target)
	if ok {
		return url, ok
	}
	url, ok = u.isHttpAlive(target)
	if ok {
		return url, ok
	}
	return url, ok
}

func (u *UrlAlive) isHttpsAlive(target string) (string, bool) {
	url := fmt.Sprintf("%s%s", "https://", target)
	ok := u.isListening(u.Client, url, "GET")
	if ok == false {
		url = ""
	}
	return url, ok
}

func (u *UrlAlive) isHttpAlive(target string) (string, bool) {
	url := fmt.Sprintf("%s%s", "http://", target)
	ok := u.isListening(u.Client, url, "GET")
	if ok == false {
		url = ""
	}
	return url, ok
}
