package domain_scan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rambleZzz/ffp/common"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ObserverWard struct {
	ObserverWardResult ObserverWardResult
}

type ObserverWardResult struct {
	Domain     string             `json:"domain"`
	Url        string             `json:"url"`
	Name       common.StringSlice `json:"name"`
	Priority   int                `json:"priority"`
	Length     int                `json:"length"`
	Title      string             `json:"title"`
	StatusCode int                `json:"status_code"`
}

func NewObserverWard() *ObserverWard {
	return &ObserverWard{}
}

func (o *ObserverWard) RunObserverWard(url string) {
	var Timeout = 30 * time.Second
	resultTempFile := common.GetTempPathFileName()
	defer os.Remove(resultTempFile)
	cmdBin := filepath.Join(common.CurrentRunPath, common.ObserverWardPath)
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "-t", url, "--fpath", filepath.Join(common.CurrentRunPath, common.WebFingerprintPath), "-j", resultTempFile)
	//Fix:指定当前路径，这样才会正确调用web_fingerprint_v3.json
	//Fix:必须指定绝对路径
	ctxt, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctxt, cmdBin, cmdArgs...)
	cmd.Dir = filepath.Join(common.CurrentRunPath, common.ObserverWardDir)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(resultTempFile)
	o.ParseObserverWardResult(resultTempFile)
}

func (o *ObserverWard) ParseObserverWardResult(resultTempFile string) {
	var obws []ObserverWardResult
	content, err := os.ReadFile(resultTempFile)
	if err != nil || len(content) == 0 {
		return
	}
	json.Unmarshal(content, &obws)
	o.ObserverWardResult = obws[0]
}
