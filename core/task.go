package core

import (
	"fmt"
	"github.com/rambleZzz/ffp/common"
	"github.com/rambleZzz/ffp/domainscan"
	"github.com/rambleZzz/ffp/ipscan"
	"github.com/urfave/cli/v2"
	"log"
	"path/filepath"
	"sync"
	"time"
)

func init() {
	common.BannerInit()
	common.CurrentRunPath, _ = common.GetCurrentPath()
	common.ResultPath = common.CreateResultDir()
	common.ReadYaml()
}

func TaskNew(cli *cli.Context) error {
	var targets []string
	//cli
	inputTarget := cli.String("u")
	inputFileName := cli.String("f")
	common.Threads = cli.Int("t")
	common.NotOutExcelFlag = cli.Bool("noe")

	if inputTarget != "" && inputFileName != "" {
		log.Println("[!]-f 与 -u 参数不可同时出现")
		log.Fatalln("[!]请输入autoScan_go -h 查看详细命令参数")
	}
	if inputTarget == "" && inputFileName == "" {
		log.Println("[!]请输入-f 或 -u 参数, -f target.txt 或 -u example.com")
		log.Fatalln("[!]请输入autoScan_go -h 查看详细命令参数")
	}
	if common.Threads == 0 {
		common.Threads = 5
	}
	if inputFileName != "" {
		inputFile := filepath.Join(common.CurrentRunPath, inputFileName)
		targets = common.RemoveDuplicate(common.GetScanTarget(inputFile))
	} else {
		if inputTarget != "" {
			targets = []string{inputTarget}
		}
	}

	// task
	common.StartTime = time.Now()
	fmt.Printf("[%v] Task started.\n", common.GetCurrentTimeStringTime())

	// urlAlive
	UrlAliveTask(targets)
	fmt.Printf("\n[%v] Domainscan scanning started, a total of (%v).\n", common.GetCurrentTimeStringTime(), len(common.AliveTargetResult))

	// domainscan
	DomainScanTask(common.AliveTargetResult)

	//无CDN IP 进行 ipscan
	common.IPResult = common.RemoveDuplicate(common.IPResult)
	fmt.Printf("\n[%v] Ipscan scanning started, a total of (%v).\n", common.GetCurrentTimeStringTime(), len(common.IPResult))
	bar := common.NewProgressbar(len(common.IPResult), "ipscan")
	for _, ip := range common.IPResult {
		if common.CheckIPV4(ip) {
			ipscan.NewIPInfo().RunIPInfo(ip)
			// 延时执行 cip.cc不可频繁反问，测试5秒正常
			time.Sleep(5 * time.Second)
			bar.Add(1)
		}
	}
	bar.Finish()
	common.PrintNewLines(1)

	// 结果输出
	common.ResultExport()
	return nil
}

// 多线程进行domainscan
func DomainScanTask(domains []string) {
	// 创建一个进度条
	bar := common.NewProgressbar(len(domains), "domainscan")
	var dwg sync.WaitGroup
	ch_domain := make(chan string, len(domains))
	//发送扫描目标到通道
	for _, domain := range domains {
		dwg.Add(1)
		ch_domain <- domain
	}
	//多线程开启扫描任务
	for i := 0; i < common.Threads; i++ {
		go func() {
			for domain := range ch_domain {
				domain_scan.NewDomainScan().Do(domain, &dwg)
				bar.Add(1)
				dwg.Done()
			}
		}()
	}
	dwg.Wait()
	bar.Finish()
	close(ch_domain)
}

// 多线程进行urlAlive
func UrlAliveTask(targets []string) {
	// 创建一个进度条
	bar := common.NewProgressbar(len(targets), "urlAlive scan")
	var uwg sync.WaitGroup
	ch_target := make(chan string, len(targets))
	//发送扫描目标到通道
	for _, target := range targets {
		uwg.Add(1)
		ch_target <- target
	}
	//多线程开启扫描任务
	for i := 0; i < common.Threads; i++ {
		go func() {
			for target := range ch_target {
				domain_scan.NewUrlAlive().Do(target, &uwg)
				// 将光标移动到终端底部
				bar.Add(1)
				uwg.Done()
			}
		}()
	}
	uwg.Wait()
	bar.Finish()
	close(ch_target)
}
