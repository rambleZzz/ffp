package domain_scan

import (
	"encoding/json"
	"fmt"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/common/customheader"
	"github.com/projectdiscovery/httpx/runner"
	"github.com/rambleZzz/ffp/common"
	"log"
	"os"
)

type Httpx struct {
	HttpxResult HttpxResult
}

func NewHttpx() *Httpx {
	return &Httpx{}
}

type HttpxResult struct {
	Domain      string             `json:"domain"`
	A           common.StringSlice `json:"a"`
	CName       common.StringSlice `json:"cname"`
	Url         string             `json:"url"`
	Host        string             `json:"host"`
	Port        string             `json:"port"`
	Title       string             `json:"title"`
	WebServer   string             `json:"webserver"`
	Tech        common.StringSlice `json:"tech"`
	ContentType string             `json:"content_type"`
	StatusCode  int                `json:"status_code"`
	Failed      bool               `json:"failed"`
	FaviconHash string             `json:"favicon_hash"`
	Scheme      string             `json:"scheme"`
	FinalUrl    string             `json:"final_url"`
}

func (h *Httpx) RunHttpx(domain string) {
	resultTempFile := common.GetTempPathFileName()
	defer os.Remove(resultTempFile)
	inputTempFile := common.GetTempPathFileName()
	defer os.Remove(inputTempFile)
	err := os.WriteFile(inputTempFile, []byte(domain), 0644)
	if err != nil {
		log.Fatal(err)
	}
	// 不输出httpx默认详情
	gologger.DefaultLogger.SetMaxLevel(levels.LevelFatal)
	options := runner.Options{
		Methods: "GET",
		//OnResult: func(r runner.Result) {
		//	// handle error
		//	fmt.Printf("\n")
		//},
		InputFile:          inputTempFile,
		CustomHeaders:      customheader.CustomHeaders{"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15063"},
		Output:             resultTempFile,
		Retries:            0,
		Threads:            50,
		Timeout:            10,
		MaxRedirects:       3,
		ExtractTitle:       true,
		StatusCode:         true,
		FollowRedirects:    true,
		JSONOutput:         true,
		NoColor:            true,
		OutputServerHeader: true,
		OutputContentType:  true,
		TechDetect:         true,
		OutputCDN:          true,
		OutputCName:        true,
		Probe:              true,
		Hashes:             "mmh3",
		//TLSGrab:            true,
	}
	if err := options.ValidateOptions(); err != nil {
		log.Fatal(err)
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		log.Fatal()
	}
	defer httpxRunner.Close()
	httpxRunner.RunEnumeration()
	h.parseHttpxResult(resultTempFile)
}

func (h *Httpx) ParseHttpxJson(content []byte) {
	resultJSON := HttpxResult{}
	err := json.Unmarshal(content, &resultJSON)
	if err != nil {
		return
	}
	// 获取host与port
	//host = resultJSON.Host
	//port, err = strconv.Atoi(resultJSON.Port)

	// 获取全部的Httpx信息
	_, err = json.Marshal(resultJSON)
	if err == nil {
		// 解析字段
		if resultJSON.Title != "" {
			h.HttpxResult.Title = resultJSON.Title
		}
		if resultJSON.WebServer != "" {
			h.HttpxResult.WebServer = resultJSON.WebServer
		}
		if resultJSON.Host != "" {
			h.HttpxResult.Host = resultJSON.Host
		}
		if resultJSON.Url != "" {
			h.HttpxResult.Url = resultJSON.Url
		}
		if resultJSON.Port != "" {
			h.HttpxResult.Port = resultJSON.Port
		}
		if resultJSON.ContentType != "" {
			h.HttpxResult.ContentType = resultJSON.ContentType
		}
		if resultJSON.Scheme != "" {
			h.HttpxResult.Scheme = resultJSON.Scheme
		}
		//if resultJSON.FaviconPath != "" {
		//	h.HttpxResult.FaviconPath = resultJSON.FaviconPath
		//}
		if resultJSON.StatusCode != 0 {
			h.HttpxResult.StatusCode = resultJSON.StatusCode
		}
		if len(resultJSON.A) != 0 {
			h.HttpxResult.A = resultJSON.A
		}
		if len(resultJSON.Tech) != 0 {
			h.HttpxResult.Tech = resultJSON.Tech
		}
		if len(resultJSON.CName) != 0 {
			h.HttpxResult.CName = resultJSON.CName
		}
		if len(resultJSON.FinalUrl) != 0 {
			h.HttpxResult.FinalUrl = resultJSON.FinalUrl
		}
		h.HttpxResult.Failed = resultJSON.Failed
	}

	fmt.Printf("\n[%v] %v\n", common.GetCurrentTimeStringTime(), h.HttpxResult)

}

// parseHttpxResult 解析httpx执行结果
func (h *Httpx) parseHttpxResult(outputTempFile string) {
	content, err := os.ReadFile(outputTempFile)
	if err != nil || len(content) == 0 {
		return
	}
	h.ParseHttpxJson(content)
}
